/**
 * getFormData serial data from form
 * @param {object} form - The jquery form object 
 * @returns {object} result - the serialized javascript object
 */
function getFormData($form) {
    var formArray = $form.serializeArray()
    var result = {}
    $.map(formArray, function(n, i) {
        result[n["name"]] = n["value"]
    })
    return result
}

/**
 * getFormData serial data from form
 * @param {object} form - The jquery form object 
 * @returns {object} result - the serialized javascript object
 */
function validateConfig(result) {
    result["port"] = parseInt(result["port"])
    if (result["port"] <= 0 || result["port"] > 65535) {
        return false
    }
    result["qps"] = parseInt(result["qps"])
    if (result["qps"] <= 0) {
        return false
    }
    result["domain_random_length"] = parseInt(result["domain_random_length"])
    if (result["domain_random_length"] <= 0) {
        return false
    }
    result["duration"] = parseInt(result["duration"])
    if (result["duration"] <= 0) {
        return false
    }
    return true
}

$(document).ready(function() {
    $(".btn-fixed-select").on("click", function() {
        $(".btn-fixed-select").removeClass("active")
        $(this).addClass("active")
    });

    $(".btn-submit").click(function() {
        // serial data
        var result = getFormData($('form[name="config"]'))
        var fixedType = $(".btn-fixed-select.active").attr("data-value") === "true" ? true : false
            // validate data
            // sende data
        result['query_type_fixed'] = fixedType
        if (validateConfig(result) === false) {
            return
        }
        $.ajax({
            type: "POST",
            url: "/start",
            data: JSON.stringify(result),
            success: function(data) {
                console.log(data)
            },
            contentType: "application/json"
        })
    })
})