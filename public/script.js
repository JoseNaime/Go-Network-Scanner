const statusElement = document.getElementById('status');
const activeIpsTable = document.getElementById('active-ips-table')
const scanDetails = document.getElementById("scan-details")

// WebSocket connection
const socket = new WebSocket('ws://localhost:5234/ws');

socket.onmessage = function(event) {
    const data = JSON.parse(event.data);
    const deviceData = {
        ip_address: data.ip_address,
        mac_address: data.mac_address
    }
    const nmapScanData = data.data.nmap_scan_data;
    console.log(data);

    const table = document.getElementById('active-ips-table'); // Ensure you have an ID for your table element
    while (table.rows.length > 1) { // Assuming the first row is the header
        table.deleteRow(1);
    }

    for (let i in nmapScanData) {
        const scanData = nmapScanData[i];
        console.log(scanData);

        let row = document.createElement("tr");
        row.classList.add("scan-table-item");

        let numberCol = document.createElement("td");
        let ipCol = document.createElement("td");
        let macCol = document.createElement("td");
        let deviceCol = document.createElement("td");

        let detailsCol = document.createElement("td");
        let detailsButton = document.createElement("button");
        detailsButton.innerText = "Details"; // Add button text
        detailsButton.onclick = () => {
            const scanDetails = document.getElementById('scan-details'); // Ensure you have an ID for your details element
            scanDetails.innerText = scanData.Details;
        };

        if (scanData.IP === deviceData.ip_address){
            scanData.MACAddress = deviceData.mac_address
            scanData.DeviceType = "THIS DEVICE"
        }

        numberCol.innerText = parseInt(i) + 1;
        ipCol.innerText = scanData.IP;
        macCol.innerText = scanData.MACAddress;
        deviceCol.innerText = scanData.DeviceType;
        detailsCol.appendChild(detailsButton);



        row.appendChild(numberCol);
        row.appendChild(ipCol);
        row.appendChild(macCol);
        row.appendChild(deviceCol);
        row.appendChild(detailsCol);

        table.appendChild(row);
    }

    console.log("Data from server", data);
    const statusElement = document.getElementById('statusElement'); // Ensure you have an ID for your status element
    statusElement.innerText = `Status: ${data.status}\nTime: ${data.time}`;
};

// Fetch initial stats
fetch('/stats')
    .then(response => response.json())
    .then(data => {
        statusElement.innerText = `Status: ${data.status}\nTime: ${data.time}`;
    });