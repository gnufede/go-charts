$(document).ready(function() {
    $(document).on('click', '.dateinput', function() {
        $(this).datepicker({
            dateFormat: "yy-mm-dd",
            showOn:'focus'
        }).focus();
    });

    $('#id_country').on('change', function () {
        $.ajax({
            url: '/city/',
            type: 'GET',
            data: {
                country: $(this).val()
            },
            success: function (data) {
                var arrayLength = data.cities.length;
                var cities = data.cities;

                $("#id_city option").remove();
                for (var i = 0; i < arrayLength; i++) {
                    city = cities[i];
                    $('#id_city').append(
                        '<option data-timezone-name="' + city.timezone + '" value=' + city.city + '>' + city.city + ' (' + city.local_timezone_name + ')</option>'
                    );
                }
                $('.selectpicker').selectpicker('refresh');
                $('#id_city').trigger('change');
            }
        });
    });

    $('#id_city').on('change', function () {
        $('#id_timezone_name').val(
            $(this).find(':selected').data('timezone-name')
        );
    });

});


Array.prototype.delete = function(value) {
    this.splice(this.indexOf(value), 1);
};
