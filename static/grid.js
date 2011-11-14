function set_hash(new_hash) {
  var href = location.href;
  var part = href.substring(0, href.length - location.hash.length);
  location.replace(part + new_hash);
}

function values_in(desc) {
  var match = desc.match(/^#(\d+)x(\d+)-(\d+)-(\d+)x(\d+)-([a-z]+)-(\d+)$/);
  return match ? match.slice(1, match.length) : false;
}

function set_values(values) {
  $('#column-width').get()[0].value = values[0];
  $('#line-height').get()[0].value = values[1];
  $('#gutter-width').get()[0].value = values[2];
  $('#cols-in-group').get()[0].value = values[3];
  $('#lines-in-group').get()[0].value = values[4];
  $('#color').get()[0].value = values[5];
  $('#column-count').get()[0].value = values[6];
}

function get_name() {
  return $('#column-width').get(0).value + 'x' +
         $('#line-height').get(0).value + '-' +
         $('#gutter-width').get(0).value + '-' +
         $('#cols-in-group').get(0).value + 'x' +
         $('#lines-in-group').get(0).value + '-' +
         $('#color').get(0).value;
}

function go() {
  var name = get_name();
  var url = '/img/grid-' + name + '.png';
  var desc = name + '-' + $('#column-count').get(0).value;
  set_hash('#' + desc);
  document.title = 'Grid ' + desc;

  var cw = parseInt($('#column-width').get(0).value);
  var gw = parseInt($('#gutter-width').get(0).value);
  var cc = parseInt($('#column-count').get(0).value);

  $('#result').get(0).value = location.protocol + '//' + location.host + url;
  $('.example')
    .css('background', 'url(' + url + ')')
    .css('width', '' + (cw * cc + gw * (cc - 1)) + 'px');

  var lh = parseInt($('#line-height').get(0).value);
  $('.example p, .example li').css('font-size', '' + (lh * 0.66) + 'px');
  $('.example h1').css('font-size', '' + (lh * 2) + 'px');

  $('.example li')
    .css('width', '' + cw + 'px')
    .css('padding-left', '' + gw + 'px');

  var total_height = 240; // TODO compute a reasonable height estimate
  $('.white, .black')
    .css('padding-left', '' + (cw + gw) + 'px')
    .css('padding-right', '' + (cw + gw) + 'px')
    .css('height', '' + total_height + 'px');


}

$(document).ready(function() {
  $('input,select').click(go).keyup(go).change(go);

  var values = values_in(location.hash);
  if (!values) {
    set_hash('#90x24-10-2x6-red-7');
    values = values_in(location.hash);
  }

  set_values(values);
  go();
});
