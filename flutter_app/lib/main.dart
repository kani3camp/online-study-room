
import 'package:flutter/material.dart';

// Import the firebase_core plugin
import 'package:firebase_core/firebase_core.dart';

import 'home_page.dart';
import 'login_page.dart';

void main() {
  runApp(App());
}

class App extends StatefulWidget {
  _AppState createState() => _AppState();
}

class _AppState extends State<App> {
  bool _initialized = false;
  bool _error = false;

  void initializeFlutterFire() async {
    try {
      await Firebase.initializeApp();
      setState(() {
        _initialized = true;
      });
    } catch(e) {
      setState(() {
        _error = true;
      });
    }
  }

  @override
  void initState() {
    initializeFlutterFire();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    if(_error) {
      return MaterialApp(
        title: 'エラー発生',
        home: Scaffold(
          body: AlertDialog(
            content: Text('エラーが発生しました'),
          ),
        ),
      );
    }

    if (!_initialized) {
      return MaterialApp(
        title: 'ローディング',
        home: Scaffold(
          body: Center(
            child: Column(
              children: [
                Text('ローディング'),
                CircularProgressIndicator()
              ]
            )
          ),
        ),
      );
    }

    return MaterialApp(
      title: 'オンライン作業部屋',
      theme: ThemeData(
        primarySwatch: Colors.blueGrey,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: Splash(),
      routes: <String, WidgetBuilder>{
        '/login': (_) => new LoginPage(),
        '/home': (_) => new MyHomePage(),
      }
    );
  }
}


class Splash extends StatefulWidget {
  @override
  _SplashState createState() => new _SplashState();
}

class _SplashState extends State<Splash> {
  @override
  void initState() {
    super.initState();
    Future.delayed(Duration(seconds: 1))
        .then((_) => handleTimeout());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Splash'),
            CircularProgressIndicator()
          ]
        ),
      ),
    );
  }

  void handleTimeout() {
    Navigator.of(context).pushReplacementNamed('/login');
  }
}
