import 'dart:convert';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/controllers/api_links.dart';
import 'package:flutter_app/pages/in_news.dart';
import 'package:http/http.dart' as http;

class NewsPage extends StatefulWidget {
  static const String pageTitle = 'お知らせ';

  @override
  _NewsPageState createState() => _NewsPageState();
}

class _NewsPageState extends State<NewsPage> {
  Future<List> _futureNewsList;

  Future<void> _init() async {
    _futureNewsList = fetchNewsList();
  }
  @override
  void initState() {
    super.initState();
    _init();
  }

  Future<bool> _onRefresh() async {
    _futureNewsList = fetchNewsList();
    await _futureNewsList;
    setState(() {});
    return true;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: RefreshIndicator(
        onRefresh: () {
          return _onRefresh();
        },
        child: FutureBuilder<List>(
          future: _futureNewsList,
          builder: (context, snapshot) {
            if (snapshot.hasData) {
              return ListView.separated(
                padding: const EdgeInsets.all(8),
                itemCount: snapshot.data.length,
                itemBuilder: (BuildContext context, int index) {
                  return Container(
                    child: ListTile(
                      title: Text(snapshot.data[index].newsBody.title),
                      subtitle: Text(
                        snapshot.data[index].newsBody.updated.year.toString() + '年'
                            + snapshot.data[index].newsBody.updated.month.toString() + '月'
                            + snapshot.data[index].newsBody.updated.day.toString() + '日'
                        ,
                        textAlign: TextAlign.right,
                      ),
                      onTap: () {
                        Navigator.of(context).pushNamed(
                            InNews.routeName,
                            arguments: InNewsArguments(snapshot.data[index].newsBody as NewsBody)
                        );
                      },
                    ),
                  );
                },
                separatorBuilder: (BuildContext context, int index) => const Divider(),
              );
            } else if (snapshot.hasError) {
              print(snapshot.error);
            } else {
            }
            return Center(child: CircularProgressIndicator());
          },
        ),
      ),
    );
  }
}

Future<List> fetchNewsList() async {
  print('fetchNewsList()');
  Map<String, String> queryParams = {
    'num_news': '10'
  };
  Uri uri = Uri.https(ApiLinks.Authority, ApiLinks.News, queryParams);
  final response = await http.get(uri);
  if (response.statusCode == 200) {
    NewsResponse roomsResponse = NewsResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
    if (roomsResponse.result == 'ok') {
      return roomsResponse.newsList;
    } else {
      throw Exception('Failed to load news list: ' + roomsResponse.message);
    }
  } else {
    throw Exception('http request failed');
  }
}

class NewsResponse {
  final String result;
  final String message;
  final List<News> newsList;

  NewsResponse({this.result, this.message, this.newsList});

  factory NewsResponse.fromJson(Map<String, dynamic> json) {
    return NewsResponse(
      result: json['result'] as String,
      message: json['message'] as String,
      newsList: (json['news_list'] as List<dynamic>).map((i) => News.fromJson(i)).toList()
    );
  }
}

class News {
  final String newsId;
  final NewsBody newsBody;

  News({this.newsId, this.newsBody});

  factory News.fromJson(Map<String, dynamic> json) {
    return News(
      newsId: json['news_id'] as String,
      newsBody: NewsBody.fromJson(json['news_body'])
    );
  }
}

class NewsBody {
  final DateTime created;
  final DateTime updated;
  final String title;
  final String textBody;

  NewsBody({this.created, this.updated, this.title, this.textBody});

  factory NewsBody.fromJson(Map<String, dynamic> json) {
    return NewsBody(
      created: DateTime.parse(json['created']).toLocal(),
      updated: DateTime.parse(json['updated']).toLocal(),
      title: json['title'] as String,
      textBody: json['text_body'] as String
    );
  }
}