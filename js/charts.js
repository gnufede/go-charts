function connect_websocket() {
    var ws = new WebSocket("ws://localhost:8080/ws?organizer=1");

    ws.onclose = function(){
        setTimeout(function(){connect_websocket()}, 5000);
    };

    ws.onmessage = function (evt) {
        var stuff = JSON.parse(evt.data)
        chart.load({json:stuff["5"]});
        channel.load({json:stuff["6"]});
        $("#tickets-sold-number").text(stuff["1"].Value[0]);
        $("#tickets-revenue-number").text(stuff["2"].Value[0]);
        $("#tickets-sold-amount").text(stuff["3"].Value[0]);
        $("#tickets-revenue-amount").text(stuff["4"].Value[0]);
//        update_value($("#tickets-sold-number"), stuff["1"].Value[0]);

    };
}

function update_value(element, new_value) {

    var old_value = parseInt(element.text());
    var new_value = parseInt(new_value);

    $({metric: old_value}).animate({metric: new_value+1}, {
        duration: 700,
        easing:'swing',
        step: function() {
            set_new_metric(element, parseInt(this.metric));
        }
    });
}

function set_new_metric($selector, metric) {
    var locale_metric = metric.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");

   // $selector.attr('data-value', metric);
    $selector.text(locale_metric);
}


var toggled = true;
var chart = c3.generate({
    color: {
        pattern: ['#495559', '#5ca648', '#A8DADC', '#C8E9A0', '#F7A278', '#413C58', '#FF7E6B', '#F9DF99', '#D6D1B1', '#B8F2E6']
    },
    data: {
        x: 'date',
        json: {},
        onclick: function (d, i) {
            if (toggled) {
                chart.toggle("children tickets");
                chart.toggle("adult tickets");
                chart.toggle("total tickets");
            }
            toggled = !toggled;
        },
        hide: ["total tickets"],
        mimeType: 'json',
        color: {
            amount: ['#495559']
        },
        axes: {
            amount: 'y2'
        },
        type: 'bar',
        types: {
            amount: 'line'
        },
        groups: [
            ['children tickets', 'adult tickets']
        ],
    },
    axis: {
        x: {
            type: 'timeseries',
            tick: {
                format: '%d-%m-%Y'
            }
        },
        y: {
            label: {
                text: 'Sold tickets',
                position: 'outer-center'
            }
        },
        y2: {
            show: true,
            tick: {
                format: function(value){
                    var format = d3.format();
                    return format(value) + '€';
                }
            },
            label: {
                text: 'Revenues',
                position: 'outer-center'
            }
        }
    },
    tooltip: {
        format: {
            value: function (value, ratio, id) {
                if(id === 'amount'){
                    var format = d3.format('');
                    return format(value) + '€';
                }else{
                    return value;
                }
            }
        }
    },
    point: {
      r: 8,
      focus: {
        expand: {
          r: 10
        }
      }
    }
});

var channel = c3.generate({
    bindto: '#donut',
    data: {
        json: {},
        //columns: [
        //    ['boxoffice', 30],
        //    ['online', 120],
        //    ['iframe', 20]
        //],
        type : 'donut',
        onclick: function (d, i) { console.log("onclick", d, i); },
        onmouseover: function (d, i) { console.log("onmouseover", d, i); },
        onmouseout: function (d, i) { console.log("onmouseout", d, i); }
    },
    donut: {
        title: "Total tickets",
        label: {
            format: function (value, ratio, id) {
                return d3.format()(value);
            }
        }
    }
});

d3.select('.container')

connect_websocket();

// setTimeout(function () {
//     chart.load({
//         columns: [
//             ['purchases', 30, 200, 100, 400, 150, 250, 350],
//             ['visits', 130, 340, 200, 500, 250, 350, 650]
//         ]
//     });
// }, 1000);

// setTimeout(function () {
//     chart.load({
//         columns: [
//             ['purchases', 30, 200, 100, 400, 150, 250, 390],
//             ['visits', 130, 340, 200, 500, 250, 350, 750 ]
//         ]
//     });
// }, 2000);


// var chart2 = c3.generate({
//     bindto: '#uv-div',
//     size: {
//         height: 150
//     },
//     bar: {
//         width: 40
//     },
//     padding: {
//         left: 100
//     },
//     color: {
//         pattern: ['#FABF62', '#ACB6DD']
//     },
//     data: {
//         x: 'x',
//         columns:
//             [
//           ['x', 'Coca Cola Music Experience', 'Arenal Sound Festival'],
//           ['purchases', 300, 400]
//           ],

//         type: 'bar',

//         color: function(inColor, data) {
//             var colors = ['#FABF62', '#ACB6DD'];
//             if(data.index !== undefined) {
//                 return colors[data.index];
//             }

//             return inColor;
//         }
//     },
//     axis: {
//         rotated: true,
//         x: {
//             type: 'category'
//         }
//     },
//     tooltip: {
//         grouped: false
//     },
//     legend: {
//         show: false
//     }
// });
