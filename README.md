unkarApp
========

Windows10で動作する5chビューアです 

**What is unkarApp**  
Golangで書いた5chビューアです  
unkarウェブサイトソースコードを元にGolang WALKを使ってWindowsアプリを作りました 

**Latest Release**  
version1.0.0.8  
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

![スクリーンショット1](https://pbs.twimg.com/media/Dbtfi3_U0AIh6Pp.jpg )  

