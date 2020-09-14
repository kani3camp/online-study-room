import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/room_page.dart';
import 'package:flutter_app/setting_page.dart';

import 'news_page.dart';

class MyHomePage extends StatefulWidget {
  MyHomePage({Key key}) : super(key: key);

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _selectedIndex = 0;
  final FirebaseAuth _auth = FirebaseAuth.instance;

  static List<Widget> _pageList = [
    RoomPage(),
    NewsPage(),
    SettingPage()
  ];

  @override
  void initState() {
    initializeFlutterFire();
    super.initState();
  }

  void initializeFlutterFire() async {
    print('initializeFlutterFire()');
    _auth.authStateChanges().listen((User user) async {
      if (user == null) {
        print('User is currently signed out!');
        // todo サインインしてない時のアプリの状態
      } else {
        print('User is signed in! : ' + user.uid);
      }
    });
  }

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    // if (!_initialized) {
    //   return Scaffold(
    //     body: Center(
    //         child: CircularProgressIndicator()
    //     ),
    //   );
    // }

    return Scaffold(
      body: _pageList[_selectedIndex],
      bottomNavigationBar: BottomNavigationBar(
        items: const <BottomNavigationBarItem>[
          BottomNavigationBarItem(
            icon: Icon(Icons.home),
            title: Text('作業'),
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.notifications),
            title: Text('お知らせ'),
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.settings),
            title: Text('設定'),
          ),
        ],
        currentIndex: _selectedIndex,
        onTap: _onItemTapped,
      ),
    );
  }
}
