import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:instachat/models/auth_user.dart';
import 'package:instachat/pages/new_chat_page.dart';
import 'package:provider/provider.dart';

import 'pages/splash_screen.dart';
import 'pages/login_page.dart';
import 'pages/chats_page.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext ctx) {
    SystemChrome.setSystemUIOverlayStyle(
      SystemUiOverlayStyle(
        statusBarColor: Colors.transparent,
      ),
    );

    return ChangeNotifierProvider(
      create: (_) => AuthUser(),
      builder: (context, _) {
        final authUser = Provider.of<AuthUser>(context);

        return MaterialApp(
          title: 'InstaChat',
          theme: ThemeData(
            brightness: Brightness.dark,
            primaryColor: Colors.black,
            accentColor: Color.fromRGBO(84, 84, 84, 1),
            scaffoldBackgroundColor: Colors.black,
            appBarTheme: AppBarTheme(color: Color.fromRGBO(36, 36, 36, 1)),
            buttonColor: Color.fromRGBO(36, 36, 36, 1),
            textTheme: TextTheme(
              headline6: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w500,
              ),
              bodyText1: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w500,
              ),
              bodyText2: TextStyle(
                fontSize: 16,
                height: 1.3,
              ),
            ),
          ),
          onGenerateRoute: (route) {
            if (authUser.account == null)
              return MaterialPageRoute(builder: (_) => SplashScreen());
            switch (route.name) {
              case LoginPage.routeName:
                return MaterialPageRoute(builder: (_) => LoginPage());
              case ChatsPage.routeName:
                return MaterialPageRoute(builder: (_) => ChatsPage());
              case NewChatPage.routeName:
                return MaterialPageRoute<bool>(builder: (_) => NewChatPage());
              default:
                return MaterialPageRoute(builder: (_) => ChatsPage());
            }
          },
        );
      },
    );
  }
}
