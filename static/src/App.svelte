<script>
  import { onMount } from 'svelte';
  import { Terminal } from '@xterm/xterm';
  import '@xterm/xterm/css/xterm.css';
  import { parse } from 'svelte/compiler';

  let terminalDiv;
  let term;
  
  let ws;
  let isConnected = false;
  
  let errorMessage = "";
  let inputBuffer = '';
  let serverMessage = "";

  let chosenAgentID = "None";
  let allAgentId = [];

  onMount(() => {
    term = new Terminal({
      cursorStyle: 'block',
      cursorWidth: 8,
      fontFamily: "'Consolas', 'Menlo', 'Courier New', monospace",
      fontSize: 14,
      letterSpacing: 0,
      lineHeight: 1,
      cols: 80,
      rows: 24,
      theme: {
        background: '#1e1e1e', // Was #ffffff
        foreground: '#d4d4d4', // Was #333
        cursor: '#d4d4d4', // Match text color
      },
      cursorBlink: true,
      allowProposedApi: true,
    });
    

    term.open(terminalDiv);
    term.write('GStrike > ');

    terminalDiv.addEventListener('mousedown', (e) => {
      term.focus();
    });

    setTimeout(() => term.focus(), 100);
    
    ws = new WebSocket("wss://localhost/ws");

    ws.onopen = () => {
      isConnected = true;
    };


    // handle server websocket messages
    ws.onmessage = (event) => {
      term.write('\r\x1b[K');  // Clear current line
      let parsed = JSON.parse(event.data);
      let type = parsed["type"];

      if (type == "beacon_callback") {
        term.write(`\r\n[*] Beacon ${parsed["agent_id"]} called home: Task (${parsed["task_id"]}) executed\r\nOutput:\r\n${parsed["output"]}\r\n\n`);
      } else if (type == "beacon_register") {
        allAgentId = [...allAgentId, parsed["id"]];
        term.write(`\r\n[*] New beacon registered: ${parsed["id"]}\r\n\n`);
      }
      term.write('GStrike > ');
      term.write(inputBuffer);
    };

    // cleanup / error
    ws.onerror = () => {
      errorMessage = "Error connecting to the WebSocket server.";
      isConnected = false;
    };
    ws.onclose = () => {
      term.writeln("\r\nWebSocket connection closed.");
      isConnected = false;
    };

    function postTaskCommand(command) {
      fetch("https://localhost:443/tasks", {
        method: "POST",
        body: JSON.stringify({
          agent_id: chosenAgentID,
          command: inputBuffer,
        }),
        headers: {
          "Content-type": "application/json; charset=UTF-8"
        }
    });
    }

    term.onData((data) => {
      // send if EOL (Enter key)
      if (data === '\r') {
        postTaskCommand(inputBuffer);
        term.write('\r\n');
        inputBuffer = '';
        term.write('GStrike > ');
      // delete char
      } else if (data === '\u007F') {
        if (inputBuffer.length > 0) {
          inputBuffer = inputBuffer.slice(0, -1);
          term.write('\b \b');
        }
        // type chars
      } else {
        inputBuffer += data;
        term.write(data);
      }
    });
  });
</script>

<main class="container">
  <div class="navbar">
    <div class="logo"><a href="/">GStrike Web UI</a></div>
    <div><a href="/">Docs</a></div>
  </div>

  <div class="layout">
    <div class="web-card">
      <div class="web-shell">
        <div bind:this={terminalDiv} class="terminal-container"></div>
        {#if errorMessage}
          <p style="color: red;">{errorMessage}</p>
        {/if}
      </div>
    </div>

    <div class="dropdown-container">
      <select class="beacon-dropdown" bind:value={chosenAgentID}>
        <option value="None">Select a Beacon</option>
        {#each allAgentId as agent}
          <option value={agent}>{agent}</option>
        {/each}
      </select>
      <div class="selected-beacon">
        Active Beacon: {chosenAgentID}
      </div>
    </div>
  
  </div>
</main>