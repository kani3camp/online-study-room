import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SettingPage extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => SettingPageState();
}

class SettingPageState extends State<SettingPage> {
  SharedPreferences _prefs;

  static const DISPLAY_NAME = 'display-name';
  static const QUICK_WORD = 'quick-word';
  static const ACCOUNT_TYPE = 'account-type';
  static const MAIL_ADDRESS = 'mail-address';
  static const SUM_STUDY_TIME = 'sum-study-time';
  static const REGISTRATION_DATE = 'registration-date';

  String _displayName;
  String _quickWord;
  String _accountType;
  String _mailAddress;
  Duration _sumStudyTime;
  DateTime _registrationDate;


  @override
  void initState() {
    _initPreferences();
    super.initState();
  }

  void _initPreferences() async {
    _prefs = await SharedPreferences.getInstance();
    setState(() {
      _displayName = _prefs.getString(DISPLAY_NAME) ?? '';
      _quickWord = _prefs.getString(QUICK_WORD) ?? '';
      _accountType = _prefs.getString(ACCOUNT_TYPE) ?? '';
      _mailAddress = _prefs.getString(MAIL_ADDRESS) ?? '';
      _sumStudyTime = new Duration(minutes: _prefs.getInt(SUM_STUDY_TIME) ?? 0);
      _registrationDate = DateTime.fromMillisecondsSinceEpoch(_prefs.getInt(REGISTRATION_DATE) ?? 0);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Center(
            child: Text('設定')
        ),
      ),
      body: ListView(
        children: [
          Container(
            child: Row(
              children: [
                Text('表示名'),
                Text(_displayName),
              ],
            ),
          ),
          Container(
            child: Row(
              children: [
                Text('ひとこと'),
                Text(_quickWord),
              ],
            ),
          ),
          Divider(),
          Container(
            child: Row(
              children: [
                Text('ログイン中のアカウント'),
                Text(_accountType),
              ],
            ),
          ),
          Container(
            child: Row(
              children: [
                Text('メールアドレス'),
                Text(_mailAddress),
              ],
            ),
          ),
          Container(
            child: Row(
              children: [
                Text('合計学習時間'),
                Text(_sumStudyTime.toString() + '分'),
              ],
            ),
          ),
          Container(
            child: Row(
              children: [
                Text('登録日'),
                Text(_registrationDate.toString()), // todo
              ],
            ),
          ),
        ],
      ),
    );
  }
}

