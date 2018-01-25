package unutil

const AffiliateMicroad_728x90 = ""
const AffiliateMicroad_300x250 = ""
const DefaultHtmlHead = `
<link rel="stylesheet" href="http://file.unkar.org/css/style.css">
`

const TimeoutMessage = `<html>
<head>
<title>503 タイムアウトしました::unkar</title>
</head>
<body>
<h1>タイムアウトしました。</h1>
<p>プログラムのバグ、もしくはアクセス過多、もしくはサーバがぶっ壊れています。</p>
<table>
<tbody>
<tr>
<td>
` + AffiliateMicroad_300x250 + `
</td>
<td>
` + AffiliateMicroad_300x250 + `
</td>
</tr>
</tbody>
</table>
</body>
</html>`

const NotImplementedMessage = `<html>
<head>
<title>501 対応していないメソッド::unkar</title>
</head>
<body>
<h1>対応していないメソッドです！</h1>
<p>色々と対応するのは面倒なのでGETとHEADリクエストのみに対応しています。</p>
<table>
<tbody>
<tr>
<td>
` + AffiliateMicroad_300x250 + `
</td>
<td>
` + AffiliateMicroad_300x250 + `
</td>
</tr>
</tbody>
</table>
</body>
</html>`

const RequestURITooLongMessage = `<html>
<head>
<title>414 無駄に長いURI::unkar</title>
</head>
<body>
<h1>無駄に長いURIです。</h1>
<p>URIが長すぎて解析するのが面倒くさいので削ってください。</p>
<table>
<tbody>
<tr>
<td>
` + AffiliateMicroad_300x250 + `
</td>
<td>
` + AffiliateMicroad_300x250 + `
</td>
</tr>
</tbody>
</table>
</body>
</html>`
