import 'package:flutter/material.dart';
import 'package:flutter_app/news_page.dart';
import 'package:flutter_app/room_page.dart';
import 'package:flutter_app/setting_page.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:google_sign_in/google_sign_in.dart';

void main() {
  runApp(App());
}

class App extends StatefulWidget {
  _AppState createState() => _AppState();
}

class _AppState extends State<App> {
  bool _initialized = false;
  bool _error = true;

  void initializeFlutterFire() async {
    try {
      // Wait for Firebase to initialize and set `_initialized` state to true
      await Firebase.initializeApp();
      setState(() {
        _initialized = true;
      });
    } catch(e) {
      // Set `_error` state to true if Firebase initialization fails
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
    // Show error message if initialization failed
    if(_error) {
      return MaterialApp(
        home: Scaffold(
          body: AlertDialog(
            title: Text('error'),
            content: Text('error'),
          ),
        ),
      );
    }

    // Show a loader until FlutterFire is initialized
    if (!_initialized) {
      return MaterialApp(
        home: Scaffold(
          body: Center(
              child: CircularProgressIndicator()
          ),
        ),
      );
    }

    return MyApp();
  }
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'オンライン作業部屋',
      theme: ThemeData(
        primarySwatch: Colors.blueGrey,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: MyHomePage(),
    );
  }
}

class MyHomePage extends StatefulWidget {
  MyHomePage({Key key}) : super(key: key);

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _selectedIndex = 0;
  bool _initialized = false;
  bool _error = false;
  final GoogleSignIn _googleSignIn = GoogleSignIn();
  final FirebaseAuth _auth = FirebaseAuth.instance;

  static List<Widget> _pageList = [
    RoomPage(),
    NewsPage(),
    SettingPage()
  ];

  @override
  void initState() {
    print('init stateします');
    initializeFlutterFire();
    super.initState();
  }

  void initializeFlutterFire() async {
    print('initializeFlutterFire()');
    try {
      setState(() async {
        _auth.authStateChanges().listen((User user) {
          if (user == null) { // todo
            print('User is currently signed out!');
          } else {
            print('User is signed in!');
          }
        });
        signInWithGoogle().then((User user) {
          print('google sign in が終わりました');
          if (user == null) {
            setState(() {
              _error = true;
            });
          }
        });

        _initialized = true;
      });
    } catch(e) {
      setState(() {
        _error = true;
      });
    }
  }

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  Future<User> signInWithGoogle() async {
    print('signInWithGoogle()');
    GoogleSignInAccount googleCurrentUser = _googleSignIn.currentUser;
    try {
      if (googleCurrentUser == null) googleCurrentUser = await _googleSignIn.signInSilently();
      if (googleCurrentUser == null) googleCurrentUser = await _googleSignIn.signIn();
      if (googleCurrentUser == null) return null;

      GoogleSignInAuthentication googleAuth = await googleCurrentUser.authentication;
      final AuthCredential credential = GoogleAuthProvider.credential(
        accessToken: googleAuth.accessToken,
        idToken: googleAuth.idToken,
      );
      final User user = (await _auth.signInWithCredential(credential)).user;

      return user;
    } catch (e) {
      print('エラー発生');
      print(e);
      return null;
    }
  }

  @override
  Widget build(BuildContext context) {
    if(_error) {
      // return todo SomethingWentWrong();
    }

    // Show a loader until FlutterFire is initialized
    if (!_initialized) {
      // todo return Loading();
    }

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
