package webssh

import (
	"bufio"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	gossh "golang.org/x/crypto/ssh"
)

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

type SshLoginModel struct {
	Addr     string 
	UserName string 
	Pwd      string 
	PemKey   string 
	PtyCols  uint32
	PtyRows  uint32
}

func sshConnect(login SshLoginModel) (client *gossh.Client, ch gossh.Channel, err error) {
	config := &gossh.ClientConfig{}
	config.SetDefaults()
	config.User = login.UserName
	if login.Pwd == "" {
		return
	} else {
		config.Auth = []gossh.AuthMethod{gossh.Password(login.Pwd)}
	}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key gossh.PublicKey) error { return nil }
	client, err = gossh.Dial("tcp", login.Addr, config)
	if err != nil {
		return
	}
	channel, incomingRequests, err := client.Conn.OpenChannel("session", nil)
	if err != nil {
		return
	}
	go func() {
		for req := range incomingRequests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()
	modes := gossh.TerminalModes{
		gossh.ECHO:          1,
		gossh.TTY_OP_ISPEED: 14400,
		gossh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, gossh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{
		Term:     "xterm",
		Columns:  login.PtyCols,
		Rows:     login.PtyRows,
		Width:    login.PtyCols * 8,
		Height:   login.PtyRows * 8,
		Modelist: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, gossh.Marshal(&req))
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("e001")
		return
	}
	ok, err = channel.SendRequest("shell", true, nil)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("e002")
		return
	}
	ch = channel
	return
}

type Request struct {
	MsgType int `json:"msg_type"`
	Token string `json:"token"`
	ServerID string `json:"server_id"`
	Command string `json:"command"`
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request, checkUserToken func(string) bool, getServerInfo func(string) SshLoginModel) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		return
	}
	isConnect := false
	var channel gossh.Channel
	var client *gossh.Client
	defer func() {
		if isConnect {
			channel.Close()
			client.Close()
		}
		ws.Close()
	}()
	done := make(chan bool, 2)
	go func() {
		defer func() {
			done <- true
		}()
		for {
			_, msgByte, err := ws.ReadMessage()
			if err != nil {
				return
			}
			req := Request{}
			err = json.Unmarshal(msgByte, &req)
			if err != nil {
				return
			}
			if req.MsgType == 1 {
				if !isConnect {
					if !checkUserToken(req.Token) {
						return
					}
					loginInfo := getServerInfo(req.ServerID)
					if loginInfo.Addr == "" {
						return
					}
					client, channel, err = sshConnect(loginInfo)
					if err != nil {
						return
					} else {
						isConnect=true
					}
				}
			} else if req.MsgType == 2 {
				if _, err := channel.Write([]byte(req.Command + "\n")); nil != err {
					return
				}
			}
		}
	}()
	go func() {
		for {
			if !isConnect {
				time.Sleep(time.Millisecond * 200)
				continue
			}
			defer func() {
				done <- true
			}()
			br := bufio.NewReader(channel)
			buf := []byte{}
			t := time.NewTimer(time.Millisecond * 100)
			defer t.Stop()
			r := make(chan rune)
			go func() {
				for {
					x, size, err := br.ReadRune()
					if err != nil {
						return
					}
					if size > 0 {
						r <- x
					}
				}
			}()
			for {
				select {
				case <-t.C:
					if len(buf) != 0 {
						err = ws.WriteMessage(websocket.TextMessage, buf)
						buf = []byte{}
						if err != nil {
							return
						}
					}
					t.Reset(time.Millisecond * 100)
				case d := <-r:
					if d != utf8.RuneError {
						p := make([]byte, utf8.RuneLen(d))
						utf8.EncodeRune(p, d)
						buf = append(buf, p...)
					} else {
						buf = append(buf, []byte("@")...)
					}
				}
			}
		}
	}()
	<-done
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
