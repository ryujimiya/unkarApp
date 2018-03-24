unkarApp
========

Windows Nativeな5ch Viewerです 

**What is unkarApp**  
GO言語で書かれたunkarウェブサイトソースコード + GO WALK によりWindowsアプリを実現しています 

**Latest Release**  
version1.0.0.0  
　[ダウンロード](https://github.com/ryujimiya/unkarApp/blob/master/publish/)  

**Important**  
　Windows10でスレッドウィンドウでJavaScriptのエラーが発生する場合  
　レジストリに下記値を追加してください。  

　HKEY_USERS  
　　S-xxxxxx  
　　　　Software  
　　　　　　Microsoft  
　　　　　　　　Internet Explorer  
　　　　　　　　　　Main  
　　　　　　　　　　　　FeatureControl  
　　　　　　　　　　　　　　FEATURE_BROWSER_EMULATION  
　　　　　　　　　　　　　　　　unkarApp.exe = (DWORD) 11001  

![スクリーンショット1](https://stat.ameba.jp/user_images/20180126/06/ryujimiya/c8/ad/j/o0586039314119092111.jpg?caw=800 )
![スクリーンショット2](https://stat.ameba.jp/user_images/20180126/06/ryujimiya/0b/a3/j/o0594063814119092145.jpg?caw=800 )
![スクリーンショット3](https://stat.ameba.jp/user_images/20180126/06/ryujimiya/11/1e/j/o0836059314119092240.jpg?caw=800 )

