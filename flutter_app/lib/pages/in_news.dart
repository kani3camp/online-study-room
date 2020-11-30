

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/pages/news_page.dart';

import '../home_page.dart';

class InNewsArguments {
  // final String newsTitle;
  // final String newsText;
  final newsBody;

  // InNewsArguments(this.newsTitle, this.newsText);
  InNewsArguments(this.newsBody);
}


class InNews extends StatefulWidget {
  static const routeName = '/in_news';

  @override
  _InNewsState createState() => _InNewsState();
}

class _InNewsState extends State<InNews> {
  NewsBody _newsBody;

  @override
  void initState() {
    super.initState();
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    final InNewsArguments arg = ModalRoute.of(context).settings.arguments;
    setState(() {
      _newsBody = arg.newsBody;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text('お知らせ 詳細',),
          leading: FlatButton(
            onPressed: () => Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName)),
            child: Icon(
              Icons.close,
              color: Colors.white,
            ),
          ),
        ),
        body: Column(children: [
          ListTile(
            title: Text(
              _newsBody.title,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 20,
              ),
            ),
          ),
          Padding(
            padding: EdgeInsets.all(5.0),
            child: Opacity(
              opacity: 0.5,
              child: Text('日時：'
                  + _newsBody.updated.year.toString()
                  + '年'
                  + _newsBody.updated.month.toString()
                  + '月'
                  + _newsBody.updated.day.toString()
                  + '日'
              ),
            ),
          ),
          Divider(),
          Align(
            child: Padding(
              padding: EdgeInsets.all(10.0),
              child: Text(_newsBody.textBody),
            ),
            alignment: Alignment.centerLeft,
          ),
        ]));
  }
}
