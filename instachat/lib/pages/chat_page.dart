import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/chat.dart';
import '../models/message.dart';
import '../services/auth_service.dart';
import '../services/chat_service.dart';
import '../widgets/insta_app_bar.dart';
import '../widgets/message.dart';
import '../widgets/message_box.dart';

class ChatPage extends StatefulWidget {
  final Chat chat;

  ChatPage(this.chat);

  @override
  _ChatPageState createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> with WidgetsBindingObserver {
  Auth auth;
  ChatService chatService;

  MessageBox _messageBox;
  ScrollController _controller;
  double _bottomInset;

  List<Message> messageCache = [];

  void updateMessages() {
    setState(() {
      messageCache = chatService.messages;
    });
  }

  void addMessage(String text) {
    final message = Message(
      senderId: chatService.auth.user.id,
      senderName: chatService.auth.user.name,
      content: text,
    );
    chatService.sendMessage(message);
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _controller.animateTo(
        _controller.position.maxScrollExtent,
        duration: Duration(milliseconds: 200),
        curve: Curves.easeOutQuad,
      );
    });
  }

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    _controller = ScrollController();
    _messageBox = MessageBox(addMessage);
  }

  @override
  void didChangeDependencies() async {
    super.didChangeDependencies();
    auth = Provider.of<Auth>(context, listen: false);

    chatService = ChatService(auth, widget.chat.id);
    chatService.connectWebsocket();
    chatService.addListener(updateMessages);
  }

  @override
  void dispose() {
    _controller.dispose();
    chatService.dispose();
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeMetrics() {
    final newBottomInset = WidgetsBinding.instance.window.viewInsets.bottom;
    if (newBottomInset != _bottomInset) {
      if (_controller.position.maxScrollExtent - _controller.offset < 10) {
        _scrollToBottom();
        _bottomInset = newBottomInset;
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: InstaAppBar(
        title: widget.chat.name,
        leading: CircleAvatar(
          backgroundImage: NetworkImage(widget.chat.imageUrl),
          radius: 18,
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              controller: _controller,
              padding: const EdgeInsets.all(12),
              itemBuilder: (_, i) {
                final message = messageCache[i];

                return message.senderId == auth.user.id
                    ? MessageRight(message: message)
                    : MessageLeft(
                        message: message,
                        isFirstMessageFromSender: false,
                      );
              },
              itemCount: messageCache.length,
            ),
          ),
          _messageBox,
        ],
      ),
    );
  }
}
