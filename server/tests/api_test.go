package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/tusharsadhwani/instachat/api"
	. "github.com/tusharsadhwani/instachat/testutils"
)

func TestMain(m *testing.M) {
	os.Setenv("GO_ENV", "TESTING")
	api.Init()

	go api.RunApp()
	m.Run()

	app := api.GetApp()
	app.Shutdown()
}

func TestHelloWorld(t *testing.T) {
	resp, err := HttpGetJson("https://localhost:5555")
	if err != nil {
		t.Fatal(err.Error())
	}
	output := string(resp)
	expected := "Hello, World 👋!"
	if output != expected {
		t.Fatalf("Expected %q, got %q", expected, output)
	}
}

func TestLogin(t *testing.T) {
	t.Run("login test user", func(t *testing.T) {
		resp, err := HttpGetJson("https://localhost:5555/test")
		if err != nil {
			t.Fatal(err.Error())
		}
		output := string(resp)
		expected := fmt.Sprintf("Welcome %s", api.TestUser.Name)
		if output != expected {
			t.Fatalf("Expected %q, got %q", expected, output)
		}
	})

	t.Run("login test user 2", func(t *testing.T) {
		resp, err := HttpGetJson("https://localhost:5555/test?testid=2")
		if err != nil {
			t.Fatal(err.Error())
		}
		output := string(resp)
		expected := fmt.Sprintf("Welcome %s", api.TestUser2.Name)
		if output != expected {
			t.Fatalf("Expected %q, got %q", expected, output)
		}
	})
}

func TestChats(t *testing.T) {
	testChat := api.Chat{
		Name:    "Test Chat",
		Address: "dbtestchat",
	}

	t.Run("create a chat and get chat by id", func(t *testing.T) {
		resp, err := HttpPostJson("https://localhost:5555/chat", testChat)
		if err != nil {
			t.Fatal(err.Error())
		}
		var respChat api.Chat
		json.Unmarshal(resp, &respChat)
		if respChat.Name != testChat.Name || respChat.Address != testChat.Address {
			t.Fatalf("Expected %#v, got %#v", testChat, respChat)
		}

		url := fmt.Sprintf("https://localhost:5555/public/chat/%d", respChat.Chatid)
		resp, err = HttpGetJson(url)
		if err != nil {
			t.Fatal(err.Error())
		}
		json.Unmarshal(resp, &respChat)
		if respChat.Name != testChat.Name || respChat.Address != testChat.Address {
			t.Fatalf("Expected %#v, got %#v", testChat, respChat)
		}
	})

	t.Run("chat id 0 test", func(t *testing.T) {
		tempChat := api.Chat{
			Address: "chatid0test",
			Name:    "Temp Chat",
		}
		_, err := HttpPostJson("https://localhost:5555/chat", tempChat)
		if err != nil {
			t.Fatal(err.Error())
		}

		_, err = HttpGetJson("https://localhost:5555/public/chat/0")
		if err == nil {
			t.Fatal("Expected error 404, got nil")
		}
		expected := "error code 404: No Chat found with id: 0"
		if err.Error() != expected {
			t.Fatalf("Expected '%v', got '%v'", expected, err)
		}

		_, err = HttpDeleteJson(fmt.Sprintf("https://localhost:5555/chat/%s", tempChat.Address))
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("delete a chat", func(t *testing.T) {
		deletionChat := api.Chat{
			Address: "deleteme",
			Name:    "Delete Me",
		}
		resp, err := HttpPostJson("https://localhost:5555/chat", deletionChat)
		if err != nil {
			t.Fatal(err.Error())
		}
		var respChat api.Chat
		json.Unmarshal(resp, &respChat)
		if respChat.Name != deletionChat.Name || respChat.Address != deletionChat.Address {
			t.Fatalf("Expected %#v, got %#v", deletionChat, respChat)
		}

		_, err = HttpDeleteJson(fmt.Sprintf("https://localhost:5555/chat/%s", deletionChat.Address))
		if err != nil {
			t.Fatal(err.Error())
		}

		_, err = HttpGetJson(fmt.Sprintf("https://localhost:5555/public/chat/%d", respChat.Chatid))
		if err == nil {
			t.Fatal("Expected error, found nil")
		}

		resp, err = HttpGetJson("https://localhost:5555/public/chat")
		if err != nil {
			t.Fatal(err.Error())
		}
		var chats []api.Chat
		json.Unmarshal(resp, &chats)
		if len(chats) != 1 {
			t.Fatalf("Expected 1 test chat to exist after deletion, found %d", len(chats))
		}
	})
}

func TestUsers(t *testing.T) {
	userOneChat := api.Chat{
		Name:    "Test User 1's Chat",
		Address: "user1chat",
	}
	userTwoChat := api.Chat{
		Name:    "Test User 2's Chat",
		Address: "user2chat",
	}

	t.Run("create chat with user 1", func(t *testing.T) {
		resp, err := HttpPostJson("https://localhost:5555/chat", userOneChat)
		if err != nil {
			t.Fatal(err.Error())
		}
		var respChat api.Chat
		json.Unmarshal(resp, &respChat)
		if respChat.Name != userOneChat.Name || respChat.Address != userOneChat.Address {
			t.Fatalf("Expected %#v, got %#v", userOneChat, respChat)
		}
		if respChat.Creatorid != api.TestUser.Userid {
			t.Fatalf("Expected creator id %d, got %d", api.TestUser.Userid, respChat.Creatorid)
		}
	})

	t.Run("create and delete chat with user 2", func(t *testing.T) {
		resp, err := HttpPostJson("https://localhost:5555/chat?testid=2", userTwoChat)
		if err != nil {
			t.Fatal(err.Error())
		}
		var respChat api.Chat
		json.Unmarshal(resp, &respChat)
		if respChat.Name != userTwoChat.Name || respChat.Address != userTwoChat.Address {
			t.Fatalf("Expected %#v, got %#v", userTwoChat, respChat)
		}
		if respChat.Creatorid != api.TestUser2.Userid {
			t.Fatalf("Expected creator id %d, got %d", api.TestUser2.Userid, respChat.Creatorid)
		}

		_, err = HttpDeleteJson(fmt.Sprintf("https://localhost:5555/chat/%s", userTwoChat.Address))
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		expected := "error code 403: 403 Forbidden"
		if err.Error() != expected {
			t.Fatalf("Expected %q ,got %q", expected, err)
		}

		_, err = HttpDeleteJson(fmt.Sprintf("https://localhost:5555/chat/%s?testid=2", userTwoChat.Address))
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestWebsockets(t *testing.T) {
	testChat := api.Chat{
		Name:    "Test Chat",
		Address: "wstestchat",
	}
	_, err := HttpPostJson("https://localhost:5555/chat", testChat)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Run("connect to websocket", func(t *testing.T) {
		url := fmt.Sprintf("https://localhost:5555/public/chat/@%s", testChat.Address)
		resp, err := HttpGetJson(url)
		if err != nil {
			t.Fatal(err.Error())
		}
		var respChat api.Chat
		json.Unmarshal(resp, &respChat)

		t.Run("send a couple messages", func(t *testing.T) {
			url := fmt.Sprintf("wss://localhost:5555/ws/%d/chat/%d", api.TestUser.Userid, respChat.Chatid)
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer conn.Close()
			defer conn.WriteMessage(websocket.CloseMessage, nil)

			testUser := api.TestUser

			msgText := "henlo"
			msg := api.WebsocketParams{
				Type: "MESSAGE",
				Message: &api.Message{
					UUID:   fmt.Sprintf("%d", rand.Uint64()),
					Chatid: &respChat.Chatid,
					Userid: &testUser.Userid,
					Text:   &msgText,
				},
			}
			_, err = WSSendAndVerify(conn, msg, testUser, respChat)
			if err != nil {
				t.Fatal(err)
			}

			msgText = "Non id pariatur dolor id Lorem ex enim proident cillum eiusmod exercitation. Laboris ut adipisicing qui minim fugiat id cupidatat velit aliquip esse commodo consequat. Excepteur deserunt duis cupidatat mollit commodo labore incididunt. Eu reprehenderit nisi commodo occaecat velit. Consequat ex officia dolor cillum exercitation incididunt occaecat ea. Culpa est veniam eiusmod aute ad adipisicing duis veniam commodo mollit exercitation dolor incididunt et."
			msg = api.WebsocketParams{
				Type: "MESSAGE",
				Message: &api.Message{
					UUID:   fmt.Sprintf("%d", rand.Uint64()),
					Chatid: &respChat.Chatid,
					Userid: &testUser.Userid,
					Text:   &msgText,
				},
			}
			_, err = WSSendAndVerify(conn, msg, testUser, respChat)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("send message with second user", func(t *testing.T) {
			url := fmt.Sprintf(
				"wss://localhost:5555/ws/%d/chat/%d?testid=2",
				api.TestUser2.Userid,
				respChat.Chatid,
			)
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer conn.Close()
			defer conn.WriteMessage(websocket.CloseMessage, nil)

			testUser := api.TestUser2
			msgText := ":D"
			msg := api.WebsocketParams{
				Type: "MESSAGE",
				Message: &api.Message{
					UUID:   fmt.Sprintf("%d", rand.Uint64()),
					Chatid: &respChat.Chatid,
					Userid: &testUser.Userid,
					Text:   &msgText,
				},
			}
			_, err = WSSendAndVerify(conn, msg, testUser, respChat)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("verify all connections receive message", func(t *testing.T) {
			testUser := api.TestUser
			msgText := "Eiusmod et veniam nulla fugiat in voluptate ullamco magna sit excepteur ex anim nulla."

			msg := api.WebsocketParams{
				Type: "MESSAGE",
				Message: &api.Message{
					UUID:   fmt.Sprintf("%d", rand.Uint64()),
					Chatid: &respChat.Chatid,
					Userid: &testUser.Userid,
					Text:   &msgText,
				},
			}

			url := fmt.Sprintf("wss://localhost:5555/ws/%d/chat/%d", api.TestUser.Userid, respChat.Chatid)
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer conn.Close()
			defer conn.WriteMessage(websocket.CloseMessage, nil)

			url2 := fmt.Sprintf("wss://localhost:5555/ws/%d/chat/%d?testid=2", api.TestUser2.Userid, respChat.Chatid)
			conn2, _, err := websocket.DefaultDialer.Dial(url2, nil)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer conn2.Close()
			defer conn2.WriteMessage(websocket.CloseMessage, nil)

			recvMsg, err := WSSendAndVerify(conn, msg, testUser, respChat)
			if err != nil {
				t.Fatal(err)
			}

			var recv2 api.WebsocketParams
			if err := conn2.ReadJSON(&recv2); err != nil {
				t.Fatal(err)
			}

			recvBytes, _ := json.Marshal(recvMsg)
			recvString := string(recvBytes)
			recv2Bytes, _ := json.Marshal(recv2)
			recv2String := string(recv2Bytes)
			if recvString != recv2String {
				t.Fatalf("expected %q, got %q", recvString, recv2String)
			}
		})

		t.Run("like a message", func(t *testing.T) {
			url = fmt.Sprintf("https://localhost:5555/public/chat/%d/message", respChat.Chatid)
			resp, err := HttpGetJson(url)
			if err != nil {
				t.Fatal(err)
			}

			var respMessagePage struct {
				Messages []api.Message
				Next     int
			}
			json.Unmarshal(resp, &respMessagePage)
			messageID := strconv.Itoa(respMessagePage.Messages[0].ID)
			msg := api.WebsocketParams{
				Type:      "Like",
				MessageID: &messageID,
			}

			url := fmt.Sprintf("wss://localhost:5555/ws/%d/chat/%d", api.TestUser.Userid, respChat.Chatid)
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer conn.Close()
			defer conn.WriteMessage(websocket.CloseMessage, nil)
			url = fmt.Sprintf("wss://localhost:5555/ws/%d/chat/%d?testid=2", api.TestUser2.Userid, respChat.Chatid)
			conn2, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer conn2.Close()
			defer conn2.WriteMessage(websocket.CloseMessage, nil)

			if err := conn.WriteJSON(msg); err != nil {
				t.Fatal(err)
			}
			var recv api.WebsocketParams
			if err := conn.ReadJSON(&recv); err != nil {
				t.Fatal(err)
			}

			msgBytes, _ := json.Marshal(msg)
			msgString := string(msgBytes)
			recvBytes, _ := json.Marshal(recv)
			recvString := string(recvBytes)
			if msgString != recvString {
				t.Fatalf("expected %q, got %q", msgString, recvString)
			}

			var recv2 api.WebsocketParams
			if err := conn2.ReadJSON(&recv2); err != nil {
				t.Fatal(err)
			}
			recv2Bytes, _ := json.Marshal(recv2)
			recv2String := string(recv2Bytes)
			if recvString != recv2String {
				t.Fatalf("expected %q, got %q", recvString, recv2String)
			}
		})
	})
}

// TODO: Reject sent messages if user not in group
// TODO: Join Group
// TODO: Message pagination
// TODO: Presigned URLs and image uploads
