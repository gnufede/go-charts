var week_or_not = false;
var stuff;

function connect_websocket() {
    var ws = new WebSocket("ws://localhost:8080/ws?organizer=1");

    ws.onclose = function(){
        setTimeout(function(){connect_websocket()}, 5000);
    };

    ws.onmessage = function (evt) {
        stuff = JSON.parse(evt.data)
        channel.load({json:stuff["6"]});

        if (week_or_not) {
            chart.load({json:stuff["5"]});
        }
        else {
            chart.load({json:stuff["7"]});
        }
        update_value($("#tickets-sold-number"), stuff["1"].Value[0]);
        update_value($("#tickets-revenue-number"),stuff["2"].Value[0]);
        update_value($("#tickets-sold-amount"),stuff["3"].Value[0]);
        update_value($("#tickets-revenue-amount"),stuff["4"].Value[0]);
        update_value($("#donut-amount"),stuff["1"].Value[0]);
    };

    $(document).ready(function () {
        $("#week-or-not").click(
            function(event) {
                event.preventDefault();
                week_or_not = !week_or_not;
                if (week_or_not) {
                    chart = generate_week(stuff["5"]);
                    $("#week-or-not").text("Change to 5 min");
                }
                else {
                    chart = generate_5m(stuff["7"]);
                    $("#week-or-not").text("Change to week");
                }
            }
        );
    });
}

function update_value(element, new_value) {

    var old_value = parseInt(element.attr('data-value')) || 0 ; //parseInt(element.text());
    var new_value = parseInt(new_value);

    $({metric: old_value}).animate({metric: new_value}, {
        duration: 499,
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


function generate_5m(json) {
return c3.generate({
    size: {
        height: 342
    },
    data: {
        x: 'minutes',
        xFormat: '%H:%M',
        json: json,
        hide: ["Total"],
        mimeType: 'json',
        colors: {
            Amount: ['#495559'],
            General: ['#71B7C4'],
            Infantil: ['#E9854A'],
            Jubilados: ['#ACCA86'],
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
    legend: {
         item: {
            onclick: function (d, i) {
                chart.toggle("General");
                chart.toggle("Gratuita");
                chart.toggle("Infantil");
                chart.toggle("Jubilados");
                chart.toggle("Total");
            },
        },
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
            },
            min: 0,
            padding: {
              bottom: 0
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
}


function generate_week (json) {
return c3.generate({
    size: {
        height: 342
    },
    data: {
        x: 'date',
        xFormat: '%Y-%m-%d',
        json: json,
        hide: ["Total"],
        mimeType: 'json',
        colors: {
            Amount: ['#495559'],
            General: ['#71B7C4'],
            Infantil: ['#E9854A'],
            Jubilados: ['#ACCA86'],
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
    legend: {
         item: {
            onclick: function (d, i) {
                chart.toggle("General");
                chart.toggle("Gratuita");
                chart.toggle("Infantil");
                chart.toggle("Jubilados");
                chart.toggle("Total");
            },
        },
    },
    axis: {
        x: {
            type: 'timeseries',
            tick: {
                format: '%d-%m-%Y',
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
}

var channel = c3.generate({
    bindto: '#donut',
    data: {
        json: {},
        type : 'donut',
        colors: {
            BoxOffice: ['#71B7C4'],
            IFrame: ['#E9854A'],
            Online: ['#ACCA86']
        },
    },
    donut: {
        label: {
            format: function (value, ratio, id) {
               return d3.format()(value);
            },
            threshold: 0

        }
    }
});

var chart = generate_5m({});

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
