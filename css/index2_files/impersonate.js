function split( val ) {
    return val.split( /,\s*/ );
}
function extractLast( term ) {
  return split( term ).pop();
}


function refresh_user_agent_token_impersonate(request, response) {
    $.ajax({
        url: "/auth/refresh_agent_token/",
        dataType: "json",
        success: function(data) {
            user_agent_token = data['access_token'];
            api_request_impersonate(request, response);
        },
        error: function (jqXHR, textStatus, errorThrown) {
            // If token invalid (removed)
            if (jqXHR.status === 412) {
                json = $.parseJSON(jqXHR.responseText);
                window.location = json['redirect_to'] + '?next=' + window.location.pathname;
            }
        }
    });
}


function api_request_impersonate(request, response) {
    $.ajax({
        url: API_URL + "/user/search",
        dataType: "json",
        data: {
            "search_term": extractLast( request.term ),
            "access_token": user_agent_token,
            "language": "es_es",
            "pagination[page]": 1,
            "pagination[per_page]": 10
        },
        error: function(jqXHR, textStatus, errorThrown) {

            // If token expired, refresh it
            if (jqXHR.status === 403) {
                json = $.parseJSON(jqXHR.responseText);

                if ('response' in json && 'code' in json['response']) {
                    // If token expired, refresh it
                    if (json['response']['code'] === 3003) {
                        refresh_user_agent_token_impersonate(request, response);
                    }
                    // If token invalid, then logout
                    else if (json['response']['code'] === 3004) {
                        window.location.href = '/logout?next=' + window.location.href;
                    }
                }
            }
        },
        success: function( data ) {
            // Remove pagination key to only leave results
            delete data["pagination"]
            response( $.map( data, function( item ) {
                return {
                    label: item.email,
                    id: item.id
                }
            }));
        }
    });
};


$(document).on("focus", "#id_impersonate_email", function (event) {

    $(this).autocomplete({
        source: function( request, response ) {
            api_request_impersonate(request, response);
        },
        select: function( event, ui ) {
            this.value = ui.item.value;
            return false;
        },
        search: function() {
            // custom minLength
            var term = extractLast( this.value );
            if ( term.length < 2 ) {
                return false;
            }
        },
        focus: function() {
            // prevent value inserted on focus
            return false;
        },
        open: function() {
            $( this ).removeClass( "ui-corner-all" ).addClass( "ui-corner-top" );
        },
        close: function() {
            $( this ).removeClass( "ui-corner-top" ).addClass( "ui-corner-all" );
        }
    });
});
