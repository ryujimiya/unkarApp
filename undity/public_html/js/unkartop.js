// うんかーTOP

(function(){

var _doc = document;
var emptyString = "";
var utop = window.utop = function(){
	this.createTable();
	this.tmpData = [];
	this.viewData = [];
	this.allReqs = 0;
	this.allReqsSec = 0;
	this.startTime = 0;
	this.allTimes = 0;
	this.rangeReqs = 0;
	this.rangeReqsSec = 0;
	this.ws = new WebSocket('ws://' + utop.conf.wshost + ':' + utop.conf.wsport + '/' + utop.conf.wspath);
	this.ws.onclose = function(){
		alert("close!");
	};
	this.ws.onmessage = jQuery.proxy(this, 'message');
};
utop.conf = {};
utop.conf.timeRange = 30;
utop.conf.timeRangeMs = utop.conf.timeRange * 1000;
utop.conf.wshost = '182.48.47.69';
utop.conf.wsport = '12345';
utop.conf.wspath = 'top';
utop.conf.size = 50;

utop.id = {};
utop.id.table = 'view';
utop.id.header = 'header';
utop.id.root = 'contents';

utop.utils = {};
utop.utils.nowTimestamp = function(){
	return +new Date();
};

utop.prototype.message = function(msg){
	var list = {};
	var obj = $.parseJSON(msg.data);
	var len = obj.length;
	var url = emptyString;
	var time_range_ms = utop.conf.timeRangeMs;
	var tmpData = this.tmpData;
	var tmpArr = this.viewData = [];
	var tmpObj = {};
	var key = emptyString;
	var i = 0;
	for(i = 0; i < len; i++){
		url = obj[i].U;
		if(list[url] === undefined){
			list[url] = obj[i];
			list[url].count = 1;
		} else {
			list[url].count++;
		}
	}
	list.timestamp = utop.utils.nowTimestamp();
	list.count = obj.length;
	tmpData.push(list);
	this.allReqs += list.count;
	if(this.startTime == 0){
		this.startTime = list.timestamp - 1000;
	}
	this.allTimes = list.timestamp - this.startTime;
	i = this.allTimes / 1000;
	i = Math.round((this.allReqs / i) * 10);
	this.allReqsSec = i / 10;
	i = 0;
	len = tmpData.length;
	while(tmpData.length > i){
		if((list.timestamp - tmpData[i].timestamp) > time_range_ms){
			i++;
		} else {
			break;
		}
	}
	if(i > 0){
		tmpData = tmpData.slice(i);
	}
	len = tmpData.length;
	for(i = 0; i < len; i++){
		var tmp = tmpData[i];
		var saki = {};
		for(key in tmp){
			if(tmpObj[key] === undefined){
				if(tmp[key].U === undefined){
					continue;
				}
				// オブジェクトをdeep cloneする
				saki = {};
				saki.count = tmp[key].count;
				saki.R = tmp[key].R;
				saki.T = tmp[key].T;
				saki.U = tmp[key].U;
				tmpObj[key] = saki;
			} else {
				tmpObj[key].count += tmp[key].count;
			}
		}
	}
	this.rangeReqs = 0;
	for(key in tmpObj){
		tmpArr.push(tmpObj[key]);
		this.rangeReqs += tmpObj[key].count;
	}
	i = Math.round((this.rangeReqs / utop.conf.timeRange) * 10);
	this.rangeReqsSec = i / 10;

	// ソートする
	tmpArr.sort(function(a, b){
		return b.count - a.count;
	});
	this.updateTable();
};

utop.prototype.updateTable = function(){
	var arr = this.viewData;
	var sec = utop.conf.timeRange;
	var tmp = 0;
	var url = emptyString;
	$('table#' + utop.id.table + ' > tbody > tr').each(function(index){
		if(arr[index] === undefined) return false;
		var tr = $('td', this);
		var it = arr[index];
		url = '/r/' + it.U;
		tmp = Math.round((it.count / sec) * 100);
		tr[0].innerHTML = it.count;
		tr[1].innerHTML = tmp / 100;
		tr[2].innerHTML = it.R;
		tr[3].innerHTML = '<a href="' + url + '">' + it.T + '</a>';
		tr[4].innerHTML = '<a href="' + url + '">' + url + '</a>';
		return true;
	});
	var jobj = $('table#' + utop.id.header + ' > tbody > tr');
	var td = $('td', jobj[0]).css('width', '50%');
	var all_times = this.allTimes / 1000;
	var day = Math.floor(all_times / 86400);
	if(day >= 1) all_times -= (day * 86400);
	var hour = Math.floor(all_times / 3600);
	if(hour >= 1) all_times -= (hour * 3600);
	var minute = Math.floor(all_times / 60);
	if(minute >= 1) all_times -= (minute * 60);
	td[1].innerHTML = 
		day + ' days, ' + 
		(hour < 10 ? '0':emptyString) + hour + ':' + 
		(minute < 10 ? '0':emptyString) + minute + ':' + 
		(all_times < 10 ? '0':emptyString) + Math.floor(all_times);

	td = $('td', jobj[1]).css('width', '50%');
	td[1].innerHTML = this.allReqs + ' reqs (  ' + this.allReqsSec + '/sec)';
	td = $('td', jobj[2]).css('width', '50%');
	td[1].innerHTML = this.rangeReqs + ' reqs (  ' + this.rangeReqsSec + '/sec)';
};

utop.prototype.createTable = function(){
	var text = [];
	text.push('<table id="' + utop.id.header + '">');
	text.push('<tbody>');
	text.push('<tr><td>utop runtime:</td><td></td></tr>');
	text.push('<tr class="wline"><td>All:</td><td></td></tr>');
	text.push('<tr class="wline"><td>R ( ' + utop.conf.timeRange + 's):</td><td></td></tr>');
	text.push('</tbody>');
	text.push('</table>');

	text.push('<br />');

	text.push('<table id="' + utop.id.table + '">');
	text.push('<thead>');
	text.push('<tr>');
	text.push('<th>REQS</th>');
	text.push('<th>REQ/S</th>');
	text.push('<th>RES</th>');
	text.push('<th>TITLE</th>');
	text.push('<th>URL</th>');
	text.push('</tr>');
	text.push('</thead>');
	text.push('<tbody>');
	for(var i = 0; i < utop.conf.size; i++){
		text.push('<tr><td></td><td></td><td></td><td></td><td></td></tr>');
	}
	text.push('</tbody>');
	text.push('</table>');
	$('#' + utop.id.root).append(text.join(emptyString));
};

$(_doc).ready(function(){
	window.$utop = new utop();
});

})();
