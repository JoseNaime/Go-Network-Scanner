const detailsElement = document.getElementById('scan-details');
const macAddressElement = document.getElementById("mac-address-text");
const ipAddressElement = document.getElementById("ip-address-text");
const maskAddressElement = document.getElementById("mask-address-text");

// WebSocket connection
const socket = new WebSocket('ws://localhost:5234/ws');

socket.onmessage = function(event) {
    const data = JSON.parse(event.data);
    const deviceData = {
        ip_address: data.data.ip_address,
        mac_address: data.data.mac_address,
        mask_address: data.data.mask_address
    }

    macAddressElement.innerText = deviceData.mac_address;
    ipAddressElement.innerText = deviceData.ip_address;
    maskAddressElement.innerText = deviceData.mask_address;

    const nmapScanData = data.data.nmap_scan_data;

    nmapScanData.sort((a,b) => { // Sort ips by last digit
        return parseInt(b.IP.split(".")[3]) - parseInt(a.IP.split(".")[3])
    })

    function createTable(nmapScanData) {
        // Calculate the maximum width for each column
        const headers = ['N', 'Ip', 'MAC', 'Device', 'Details'];
        const footer = ["", data.status, formatDate(data.time), "" , ""]
        const maxLengths = headers.map((header, i) => Math.max(footer[i].length, header.length));
        console.log(maxLengths)

        nmapScanData.forEach((row, index) => {
            console.log(row.MACAddress, row.DeviceType)
            if (row.MACAddress === "") row.MACAddress = deviceData.mac_address;
            if (row.DeviceType === "") row.DeviceType = "This Device";

            const values = [String(index + 1),
                row.IP,
                row.MACAddress,
                row.DeviceType,
                'Button'];
            values.forEach((value, i) => {
                maxLengths[i] = Math.max(maxLengths[i], value.length);
            });
        });

        const createSeparator = () => {
            return '+-' + maxLengths.map(len => '-'.repeat(len)).join('-+-') + '-+';
        };

        const createRow = (values, isHeader = false, isFooter = false, index = null) => {
            const cells = values.map((value, i) => value.padEnd(maxLengths[i]));
            if (isHeader) {
                return `<span class="header">| ${cells.join(' | ')} |</span>`;
            } else if (isFooter){
                return `<span class="header">| ${cells.join('   ')} |</span>`;
            }
            const rowHTML = cells.map((cell, i) => {
                if (i === cells.length - 1) {
                    return `<button class="details-button" data-index="${index}">Details</button>`.padEnd(maxLengths[i]);
                }
                return cell;
            });
            return `<span class="cell">| ${rowHTML.join(' | ')} |</span>`;
        };

        // Build the table
        let table = `<span class="border">${createSeparator()}</span>\n`;
        table += createRow(headers, true) + '\n';
        table += `<span class="border">${createSeparator()}</span>\n`;

        nmapScanData.forEach((row, index) => {
            const values = [String(index + 1), row.IP, row.MACAddress, row.DeviceType, 'Details'];
            table += createRow(values, false, false, index) + '\n';
        });

        table += `<span class="border">${createSeparator()}</span>\n`;
        table += createRow(footer, false, true) + '\n';
        table += `<span class="border">${createSeparator()}</span>`;
        return table;
    }

    // Generate the table and insert into the HTML
    document.getElementById('asciiTable').innerHTML = createTable(nmapScanData);

    // Add event listeners to the buttons
    document.querySelectorAll('.details-button').forEach(button => {
        button.addEventListener('click', function() {
            const index = this.getAttribute('data-index');
            const rowData = nmapScanData[index];
            detailsElement.innerText =  rowData.Details;
        })}
    )

    console.log("Data from server", data);
};

// Fetch initial stats
fetch('/stats')
    .then(response => response.json())
    .then(data => {
        statusElement.innerText = `Status: ${data.status}\nTime: ${data.time}`;
    });

function formatDate(dateStr){
    const date = new Date(dateStr);

    const pad = (num) => num.toString().padStart(2, '0');

    const year = date.getFullYear();
    const month = pad(date.getMonth() + 1); // Months are zero-based
    const day = pad(date.getDate());

    const hours = pad(date.getHours());
    const minutes = pad(date.getMinutes());
    const seconds = pad(date.getSeconds());

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}