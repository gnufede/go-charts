
var toggled = true;
var chart = c3.generate({
    color: {
        pattern: ['#495559', '#5ca648', '#A8DADC', '#C8E9A0', '#F7A278', '#413C58', '#FF7E6B', '#F9DF99', '#D6D1B1', '#B8F2E6']
    },
    data: {
        x: 'date',
        // url: 'http://private-dee9c-fedemon.apiary-mock.com/1.2/json/graph2.json',
        json: {
            "date": ["2013-01-01", "2013-01-02", "2013-01-03", "2013-01-04", "2013-01-05", "2013-01-06", "2013-01-07"],
            "amount": [1500, 1000, 3000, 4000, 0, 2500, 3000],
            "children tickets":[ 10, 50, 25, 100, 0, 0, 100],
            "adult tickets":[ 20, 150, 75, 300, 150, 250, 250],
            "total tickets":[ 30, 200, 100, 400, 150, 250, 350]
        },
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
        columns: [
            ['boxoffice', 30],
            ['online', 120],
            ['iframe', 20]
        ],
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