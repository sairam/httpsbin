// Adds to the list of recent bins
function AddToRecent(binID, name, expiresAt, count) {
  if (name == undefined || name == "") {
    name = binID;
  }
  var bin = {
    "binID": binID,
    "name": name,
    "expiresAt": expiresAt,
    "accessedAt": Math.floor(Date.now() / 1000),
    "count": count
  }

  var json = JSON.parse(localStorage[BinKey]);
  json[binID] = bin;
  localStorage[BinKey] = JSON.stringify(json);
}
function DeleteExpiredBins() {

  // delete json["binID"]
}

function RenderRecentBins(id) {
  var json = JSON.parse(localStorage[BinKey]);
  var arr = [];
  var keylist = Object.keys(json);
  for ( i=0 ; i < keylist.length; i++ ) {
    arr[i] = json[keylist[i]]
  }

  arr.sort(function(a,b) {
    return b.accessedAt - a.accessedAt;
  });

  var data = "";
  for( i=0 ; i < arr.length; i++ ) {
    bin = arr[i];
    var d = '<li class="list-group-item"><a href="/'+bin.binID+'?inspect">'+bin.name+'</a> ('+ bin.count + ')</li>' ;
    data += d;
  }
  if (arr.length == 0) {
    data = '<li class="list-group-item text-center"><a class="btn btn-primary" onclick="document.forms[0].submit();">Create a new Bin</a></li>';
  }
  var bins = document.getElementById(id);
  bins.innerHTML = data;
}

var BinKey = "recentBins"
if (localStorage[BinKey] == undefined) {
  localStorage[BinKey] = "{}";
}

function CreateNewBin() {
  document.forms[0].submit();
}
