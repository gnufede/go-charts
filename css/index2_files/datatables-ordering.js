// Plugin for ordering formatted numbers in Spanish
jQuery.extend( jQuery.fn.dataTableExt.oSort, {
    "formatted-num-pre": function ( a ) {
        a = (a==="-") ? 0 : a.replace( ".", "" );
        a = (a==="-") ? 0 : a.replace( ",", "." );
        return parseFloat( a );
   },

    "formatted-num-asc": function ( a, b ) {
        return a - b;
    },

    "formatted-num-desc": function ( a, b ) {
        return b - a;
    }
});


/*
 * Custom version of Ronan Guilloux date order plugin
 * supports: dd/mm/YYYY hh:ii:ss and dd/mm/YY (with and without time)
*/
jQuery.extend( jQuery.fn.dataTableExt.oSort, {
    "date-eu-pre": function ( a ) {
        if ($.trim(a) != '') {
            var date_array = $.trim(a).split(' ');

            frTimea = null;
            if (date_array.length === 2) {
                var frTimea = date_array[1].split(':');
            }
            var frDatea = date_array[0].split('/');

            if (frTimea === null) {
                var x = (frDatea[2] + frDatea[1] + frDatea[0]) * 1;
            } else {
                var x = (frDatea[2] + frDatea[1] + frDatea[0] + frTimea[0] + frTimea[1] + frTimea[2]) * 1;
            }
        } else {
            var x = 10000000000000; // = l'an 1000 ...
        }

        return x;
    },

    "date-eu-asc": function ( a, b ) {
        return a - b;
    },

    "date-eu-desc": function ( a, b ) {
        return b - a;
    }
});
