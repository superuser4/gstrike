<script>
  let ws;
  let isConnected = false;
  let myMsg = "";
  let srvMsg = "";
  let errorMessage = "";

  import { onMount } from 'svelte';

  onMount(() => {
    ws = new WebSocket("wss://localhost/ws");

    ws.onopen = () => {
      console.log("WebSocket connection opened");
      isConnected = true;
    };

    ws.onmessage = (event) => {
      srvMsg = event.data;
    };

    ws.onerror = (error) => {
      errorMessage = "Error connecting to the WebSocket server.";
      isConnected = false;
    };

    ws.onclose = () => {
      isConnected = false;
    };
  });

  const sendMessage = () => {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(myMsg);
    } else {
      errorMessage = "Unable to send message, WebSocket not open.";
    }
  };

  const handleKeydown = (event) => {
    if (event.key === "Enter") {
      sendMessage();
    }
  };
</script>

<main class="container">
  <div class="navbar">
    <div class="logo">Gstrike Web UI</div>
    <div>
      <a href="/">Home</a>
    </div>
  </div>

  <div class="web-shell">

    <input 
      bind:value={myMsg}
      placeholder="send commands to beacon" 
      on:keydown={handleKeydown} 
    />

    {#if srvMsg}
      <p><strong>Response from Beacon:</strong> {srvMsg}</p>
    {/if}

    {#if errorMessage}
      <p style="color: red;">{errorMessage}</p>
    {/if}
  </div>
</main>
