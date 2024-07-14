const statusElement = document.getElementById('status');
const activeIpsTable = document.getElementById('active-ips-table')
const scanDetails = document.getElementById("scan-details")

// WebSocket connection
const socket = new WebSocket('ws://localhost:5234/ws');

socket.onmessage = function(event) {
    const data = JSON.parse(event.data);
    const nmapScanData = data.data.nmap_scan_data
    console.log(data)

    for (let i in nmapScanData){
        const scanData = nmapScanData[i]
        console.log(scanData)
        let row = document.createElement("tr")
        let numberCol = document.createElement("td")
        let ipCol = document.createElement("td")
        let macCol = document.createElement("td")
        let deviceCol = document.createElement("td")

        let detailsCol = document.createElement("td")
        let detailsButton = document.createElement("button")
        detailsButton.onclick = () => scanDetails.innerText = scanData.Details

        numberCol.innerText = i
        ipCol.innerText = scanData.IP
        macCol.innerText = scanData.MACAddress
        deviceCol.innerText = scanData.DeviceType
        detailsCol.appendChild(detailsButton)

        row.appendChild(numberCol)
        row.appendChild(ipCol)
        row.appendChild(macCol)
        row.appendChild(deviceCol)
        row.appendChild(detailsCol)

        activeIpsTable.appendChild(row)
    }
    console.log("Data from server", data)
    statusElement.innerText = `Status: ${data.status}\nTime: ${data.time}`;
};

// Fetch initial stats
fetch('/stats')
    .then(response => response.json())
    .then(data => {
        statusElement.innerText = `Status: ${data.status}\nTime: ${data.time}`;
    });