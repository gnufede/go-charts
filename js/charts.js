function connect_websocket() {
    var ws = new WebSocket("ws://localhost:8080/ws?organizer=1");

    ws.onclose = function(){
        setTimeout(function(){connect_websocket()}, 5000);
    };

    ws.onmessage = function (evt) {
        var stuff = JSON.parse(evt.data)
        chart.load({json:stuff["7"]});
        channel.load({json:stuff["6"]});
        update_value($("#tickets-sold-number"), stuff["1"].Value[0]);
        update_value($("#tickets-revenue-number"),stuff["2"].Value[0]);
        update_value($("#tickets-sold-amount"),stuff["3"].Value[0]);
        update_value($("#tickets-revenue-amount"),stuff["4"].Value[0]);
        update_value($("#donut-amount"),stuff["1"].Value[0]);

    };
}

function update_value(element, new_value) {

    var old_value = parseInt(element.attr('data-value')) || 0 ; //parseInt(element.text());
    var new_value = parseInt(new_value);

    $({metric: old_value}).animate({metric: new_value}, {
        duration: 700,
        easing:'swing',
        step: function() {
            set_new_metric(element, parseInt(this.metric));
        },
        complete: function() {
            set_new_metric(element, new_value);
        }

    });
}

function set_new_metric($selector, metric) {
    var locale_metric = metric.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");

    $selector.attr('data-value', metric);
    $selector.text(locale_metric);
}


var toggled = true;
var chart = c3.generate({
    // color: {
    //     pattern: ['#495559', '#5ca648', '#A8DADC', '#C8E9A0', '#F7A278', '#413C58', '#FF7E6B', '#F9DF99', '#D6D1B1', '#B8F2E6']
    // },
    size: {
        height: 342
    },
    data: {
        x: 'minutes',
        xFormat: '%H:%M',
        json: {},
        onclick: function (d, i) {
            if (toggled) {
                chart.toggle("General");
                chart.toggle("Gratuita");
                chart.toggle("Infantil");
                chart.toggle("Jubilados");
                chart.toggle("Total");
            }
            toggled = !toggled;
        },
        hide: ["Total"],
        mimeType: 'json',
        color: {
            Amount: ['#495559'],
            General: ['#71B7C4'],
            Infantil: ['#E9854A'],
            Jubilados: ['#E9854A'],
            Gratuita: ['#FBDD8D'],
            Total: ['#71B7C4']
        },
        axes: {
            Amount: 'y2'
        },
        type: 'bar',
        types: {
            Amount: 'line'
        },
        groups: [
            ['General', 'Infantil', 'Jubilados', 'Gratuita']
        ],
    },
    axis: {
        x: {
            type: 'timeseries',
            tick: {
                format: '%H:%M'
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
                if(id === 'Amount'){
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
        type : 'donut'
    },
    donut: {
        label: {
            format: function (value, ratio, id) {
                return d3.format()(value);
            }
        }
    }
});

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
