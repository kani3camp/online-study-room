import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/pages/room_page.dart';
import 'package:flutter_app/pages/setting_page.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';
import 'package:url_launcher/url_launcher.dart';

import 'login_page.dart';
import 'main.dart';
import 'pages/news_page.dart';

class MyHomePage extends StatefulWidget {
  static const routeName = '/home';

  MyHomePage({Key key}) : super(key: key);

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _selectedIndex = 0;

  List<Page> _pageList = [
    Page(RoomPage(), RoomPage.pageTitle),
    Page(NewsPage(), NewsPage.pageTitle),
    Page(SettingPage(), SettingPage.pageTitle),
  ];

  @override
  void initState() {
    initializeFlutterFire();
    super.initState();
  }

  void initializeFlutterFire() async {
    print('initializeFlutterFire()');
    FirebaseAuth.instance.authStateChanges().listen((User user) async {
      if (user == null) {
        print('User is currently signed out!');
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

  Drawer _drawer = Drawer(
    child: ListView(
      children: [
        DrawerHeader(
          child: Center(
            child: Text(
              'メニュー',
              style: TextStyle(
                fontSize: 24,
              ),
            ),
          ),
        ),
        ListTile(
          leading: Icon(MdiIcons.fileDocument),
          title: Text(
            '利用規約'
          ),
          onTap: () async {
            await launch('https://online-study-space.web.app/terms_of_service');
          },
          trailing: Icon(Icons.launch),
        ),
        ListTile(
          leading: Icon(Icons.privacy_tip),
          title: Text(
            'プライバシーポリシー'
          ),
          onTap: () async {
            await launch('https://online-study-space.web.app/privacy_policy');
          },
          trailing: Icon(Icons.launch),
        ),
        ListTile(
          leading: Icon(MdiIcons.youtube),
          title: Text(
            'YouTubeチャンネル'
          ),
          onTap: () async {
            final url = 'https://www.youtube.com/channel/UCXuD2XmPTdpVy7zmwbFVZWg';
            await launch(url);
          },
          trailing: Icon(Icons.launch),
        ),
        ListTile(
          leading: Icon(MdiIcons.twitter),
          title: Text(
            'Twitter'
          ),
          onTap: () async {
            final url = 'https://twitter.com/osr_soraride';
            await launch(url);
          },
          trailing: Icon(Icons.launch),
        ),
        ListTile(
          leading: Icon(Icons.email),
          title: Text(
            'お問い合わせ'
          ),
          onTap: () async {
            final url = 'https://online-study-space.web.app/contact_form';
            await launch(url);
          },
          trailing: Icon(Icons.launch),
        ),
        ListTile(
          leading: Icon(Icons.account_circle),
          title: Text(
            '開発者'
          ),
          onTap: () async {
            final url = 'https://twitter.com/sorarideblog';
            await launch(url);
          },
          trailing: Icon(Icons.launch),
        ),
      ],
    ),
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.white,
        centerTitle: true,
        title: Text(
          _pageList[_selectedIndex].title,
          style: TextStyle(
            color: Theme.of(context).primaryColor
          ),
        ),
      ),
      body: _pageList[_selectedIndex].page,
      drawer: _drawer,
      bottomNavigationBar: BottomNavigationBar(
        items: const <BottomNavigationBarItem>[
          BottomNavigationBarItem(
            icon: Icon(Icons.home),
            label: '作業',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.notifications),
            label: 'お知らせ',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.settings),
            label: '設定',
          ),
        ],
        currentIndex: _selectedIndex,
        onTap: _onItemTapped,
      ),
    );
  }
}


class Page {
  final page;
  final title;

  Page(this.page, this.title);
}