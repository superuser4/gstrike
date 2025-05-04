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
        term.write(`\r\n${event.data}\r\n\n`);
        term.write('GStrike > ');
        term.write(inputBuffer);  // Restore current input
      } else if (type == "beacon_register") {
        allAgentId = [...allAgentId, parsed["agentID"]];
      }
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

    function taskCommand(command) {

    }

    term.onData((data) => {
      // send if EOL (Enter key)
      if (data === '\r') {
        taskCommand(inputBuffer);
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
    <div class="logo">Gstrike Web UI</div>
    <div><a href="/">Home</a></div>
  </div>

  <div class="web-shell">
    <div bind:this={terminalDiv} class="terminal-container"></div>
    {#if errorMessage}
      <p style="color: red;">{errorMessage}</p>
    {/if}
  </div>

  <div class="left-card-list">
    <h3>
      Beacon list
    </h3>
    <h4>Selected Beacon: {chosenAgentID}</h4>
    <ul>  
      {#each allAgentId as agent}
        <li>
          <button on:click={chosenAgentID = agent}>
            {agent}
          </button>
        </li>
      {/each}
    </ul>
  </div>
</main>