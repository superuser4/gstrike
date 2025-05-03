<script>
  import { onMount } from 'svelte';
  import { Terminal } from '@xterm/xterm';
  import '@xterm/xterm/css/xterm.css';

  let terminalDiv;
  let term;
  let ws;
  let isConnected = false;
  let errorMessage = "";
  let inputBuffer = '';

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
    term.write('$ ');

    ws = new WebSocket("wss://localhost/ws");

    ws.onopen = () => {
      isConnected = true;
    };

    ws.onmessage = (event) => {
      term.write('\r\x1b[K');  // Clear current line
      term.write(`\r\n${event.data}\r\n`);
      term.write('$ ');
      term.write(inputBuffer);  // Restore current input
    };

    ws.onerror = () => {
      errorMessage = "Error connecting to the WebSocket server.";
      isConnected = false;
    };

    ws.onclose = () => {
      term.writeln("\r\nWebSocket connection closed.");
      isConnected = false;
    };

    term.onData((data) => {
      if (data === '\r') {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(inputBuffer);
          term.write('\r\n');
        }
        inputBuffer = '';
        term.write('$ ');
      } else if (data === '\u007F') {
        if (inputBuffer.length > 0) {
          inputBuffer = inputBuffer.slice(0, -1);
          term.write('\b \b');
        }
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
</main>