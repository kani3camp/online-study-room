
import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/controllers/shared_preferences.dart';
import 'package:flutter_app/home_page.dart';
import 'package:google_sign_in/google_sign_in.dart';

class LoginPage extends StatefulWidget {
  static const routeName = '/login';
  LoginPage({Key key}) : super(key: key); // what is this ?

  @override
  _LoginPageState createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {

  final GoogleSignIn _googleSignIn = GoogleSignIn();
  final FirebaseAuth _auth = FirebaseAuth.instance;

  @override
  void initState() {
    signInWithGoogle();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Text('ログイン中...'),
      )
    );
  }

  Future<void> signInWithGoogle() async {
    print('signInWithGoogle()');
    GoogleSignInAccount googleCurrentUser = _googleSignIn.currentUser;
    print('p0');
    try {
      // if (googleCurrentUser == null) googleCurrentUser = await _googleSignIn.signInSilently();
      if (googleCurrentUser == null) googleCurrentUser = await _googleSignIn.signIn();
      if (googleCurrentUser == null) return null;

      GoogleSignInAuthentication googleAuth = await googleCurrentUser.authentication;
      final AuthCredential credential = GoogleAuthProvider.credential(
        accessToken: googleAuth.accessToken,
        idToken: googleAuth.idToken,
      );
      final User user = (await _auth.signInWithCredential(credential)).user;

      // final SharedPrefs _prefs = await SharedPrefs.create();
      SharedPrefs _prefs = new SharedPrefs();
      await _prefs.init();
      // await _prefs.setIdToken(await user.getIdToken());
      await _prefs.setUserId(user.uid);
      await _prefs.setMailAddress(user.email);
      await _prefs.setAccountType(user.providerData[0].providerId);

      print('google sign in が終わりました');
      if (user == null) {
        // setState(() {
        //   _error = true;
        // }); todo
      } else {
        Navigator.of(context).pushReplacementNamed(MyHomePage.routeName);
      }
    } catch (e) {
      print('エラー発生');
      print(e);
      return null;
    }
  }
}