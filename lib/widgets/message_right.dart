import 'package:flutter/material.dart';

import '../models/message.dart';
import '../widgets/likeable.dart';

class MessageRight extends StatelessWidget {
  final Message message;

  const MessageRight({Key key, this.message}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Align(
      alignment: Alignment.centerRight,
      child: Padding(
        padding: const EdgeInsets.all(3),
        child: Builder(builder: (_) {
          return Likeable(
            child: Container(
              decoration: BoxDecoration(
                color: Theme.of(context).appBarTheme.color,
                borderRadius: BorderRadius.circular(24),
              ),
              child: Padding(
                padding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                child: Text(message.content),
              ),
            ),
          );
        }),
      ),
    );
  }
}
