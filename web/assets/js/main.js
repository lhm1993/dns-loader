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

    if (result["server"] == "") {
        toastr.error('DNS server is empty', 'Config Error')
        return false
    }

    if (result["domain"] == "") {
        toastr.error('query domain is empty', 'Config Error')
        return false
    }
    if (result["query_type"] == "") {
        result["query_type"] = "A"
    }
    result["port"] = isNaN(parseInt(result["port"])) ? 53 : parseInt(result["port"])
    if (result["port"] <= 0 || result["port"] > 65535) {
        toastr.error('Port number should be in [0-65535]', 'Port Error')
        return false
    }

    result["qps"] = isNaN(parseInt(result["qps"])) ? 100 : parseInt(result["qps"])
    if (result["qps"] <= 0) {
        toastr.error('QPS number should be larger than 0', 'QPS Error')
        return false
    }
    result["domain_random_length"] = isNaN(parseInt(result["domain_random_length"])) ? 5 : parseInt(result["domain_random_length"])
    if (result["domain_random_length"] <= 0) {
        toastr.error('Random length number should be larger than 0', 'Length Error')
        return false
    }
    result["duration"] = isNaN(parseInt(result["duration"])) ? 60 : parseInt(result["duration"])
    if (result["duration"] <= 0) {
        toastr.error('Duration time should be larger than 0', 'Duration Time Error')
        return false
    }
    return true
}

$(document).ready(function() {
    $(".btn-fixed-select").on("click", function() {
        $(".btn-fixed-select").removeClass("active")
        $(this).addClass("active")
    });

    $(".config-submit").click(function() {
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
    $(".small-delete-button").click(function() {
        var ipWithPort = $(this).attr("data-item")
        var data = {
            "ipaddress": ipWithPort.split(":")[0],
            "port": parseInt(ipWithPort.split(":")[1])
        }
        $.ajax({
            type: "DELETE",
            url: "/nodes",
            data: JSON.stringify(data),
            success: function(data) {
                toastr.info("delete success")
                window.location.reload()
            },
            error: function(err) {
                if (err && err.responseJSON && err.responseJSON.message) {
                    toastr.error("Delete fail", err.responseJSON.message)
                } else {
                    toastr.error("Delete fail")
                }
            },
            contentType: "application/json"
        })
    })

    $(".small-ping-button").click(function() {
        var ipWithPort = $(this).attr("data-item")
        var data = {
            "ipaddress": ipWithPort.split(":")[0],
            "port": parseInt(ipWithPort.split(":")[1])
        }
        $.ajax({
            type: "POST",
            url: "/ping",
            data: JSON.stringify(data),
            success: function(data) {
                toastr.success("ping success")
            },
            error: function(err) {
                if (err && err.responseJSON && err.responseJSON.message) {
                    toastr.error("ping fail", err.responseJSON.message)
                } else {
                    toastr.error("ping fail")
                        // change the color of status to black
                }
            },
            contentType: "application/json"
        })
    })


    $('.config-kill').click(function() {
        console.log("stop signal send to master server")
        $.ajax({
            type: "GET",
            url: "/stop",
            success: function(response) {
                toastr.success("stop traffic success")
            },
            error: function(err) {
                if (err && err.responseJSON && err.responseJSON.message) {
                    toastr.error("Error", err.responseJSON.message)
                } else {
                    toastr.error("Error", "ServerFail")
                }
            },
            contentType: "application/json"
        })
    })

    $('.new-agent').click(function() {
        var data = getFormData($("form[name='new-agent']"))
        if (typeof data.ipaddress === 'undefined' || data.ipaddress === "") {
            toastr.error('IP address does not exist', 'IP Error')
            return
        }
        data.port = parseInt(data.port)
        if (data.port < 0) {
            toastr.error('IP address does not exist', 'IP Error')
            return
        }
        $(".add-node-loading").removeClass("hide")
        $.ajax({
            type: "POST",
            url: "/nodes",
            data: JSON.stringify(data),
            success: function(response) {
                $(".add-node-loading").addClass("hide")
                    // console.log(data)
                    //     // Add to list 
                window.location.reload()
            },
            error: function(err) {
                $(".add-node-loading").addClass("hide")
                if (err && err.responseJSON && err.responseJSON.message) {
                    toastr.error("Add new node fail", err.responseJSON.message)
                } else {
                    toastr.error("Add new node fail", "ServerFail")
                }
            },
            contentType: "application/json"
        })
    })
})