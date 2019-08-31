window.onload = function () {

    const inputFieldJob = document.getElementById("inputFieldJob")
    const inputFieldLocation = document.getElementById("inputFieldLocation")
    const submitButton = document.getElementById("submitButton")


    var options = {
        strings: ["p r i n t e r n."],
        typeSpeed: 150,
        showCursor: false,
    }

    var typed = new Typed(".main-title", options);

    submitButton.onclick = async function (e) {

        let body = {
            jobTitle: inputFieldJob.value,
            jobLocation: inputFieldLocation.value
        }


        let response = await fetch("/printer", {
            method: "POST",
            headers: {
                "Accept": "application/json"
            },
            body: JSON.stringify(body)
        })

        console.log(response.json(4))
    }


}