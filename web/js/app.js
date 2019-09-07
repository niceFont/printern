window.onload = function () {

    const inputFieldJob = document.getElementById("inputFieldJob")
    const inputFieldLocation = document.getElementById("inputFieldLocation")
    const submitButton = document.getElementById("submitButton")


    var options = {
        strings: ["p r i n t e r n."],
        typeSpeed: 130,
        showCursor: false,
    }

    var typed = new Typed(".main-title", options);

    window.onscroll = function () {
        console.log("wow")
    }
    submitButton.onclick = function (e) {

        let body = {
            jobTitle: inputFieldJob.value,
            jobLocation: inputFieldLocation.value
        }


        fetch("/printer", {
            method: "POST",
            headers: {
                "Accept": "application/json"
            },
            body: JSON.stringify(body)
        })
            .then(response => response.json())
            .then(data => {
                let arr = []

                for (let t in data) {
                    arr.push({ [t]: data[t] })
                }
                return arr.sort((a, b) => Object.values(b)[0] - Object.values(a)[0])
            })
            .then(sortedData => {
                let frag = document.createDocumentFragment()
                let resultNode = document.getElementById("result")

                let node = document.getElementById("response")

                if (!node) {

                    node = document.createElement("canvas")
                    node.id = "response"

                    frag.appendChild(node)
                    resultNode.appendChild(frag)
                }


                let ctx = node.getContext("2d")
                ctx.clearRect(0, 0, ctx.width, ctx.height)
                let chart = new Chart(ctx, {
                    type: "bar",
                    data: {
                        labels: sortedData.map(x => Object.keys(x)[0]),
                        datasets: [{
                            label: "Technology Frequency",
                            data: sortedData.map(x => Object.values(x)[0]),
                            backgroundColor: [
                                'rgba(255, 99, 132, 0.2)',
                                'rgba(54, 162, 235, 0.2)',
                                'rgba(255, 206, 86, 0.2)',
                                'rgba(75, 192, 192, 0.2)',
                                'rgba(153, 102, 255, 0.2)',
                                'rgba(255, 159, 64, 0.2)',
                                'rgba(255, 99, 132, 0.2)',
                                'rgba(255, 99, 132, 0.2)',
                                'rgba(255, 99, 132, 0.2)',
                            ],
                            borderColor: [
                                'rgba(255, 99, 132, 1)',
                                'rgba(255, 99, 132, 1)',
                                'rgba(255, 99, 132, 1)',
                                'rgba(255, 99, 132, 1)',
                                'rgba(54, 162, 235, 1)',
                                'rgba(255, 206, 86, 1)',
                                'rgba(75, 192, 192, 1)',
                                'rgba(153, 102, 255, 1)',
                                'rgba(255, 159, 64, 1)'
                            ],
                            borderWidth: 1
                        }]
                    },
                    options: {
                        scales: {
                            yAxes: [{
                                ticks: { beginAtZero: true }
                            }]
                        }
                    }
                })

                window.scrollBy({ top: document.getElementById("response").offsetTop, left: 0, behaviour: "smooth" })
            })

    }


}