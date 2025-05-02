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
    <div class="logo">Gobricked Web UI</div>
    <div>
      <a href="/">Home</a>
    </div>
  </div>

  <div class="main-content">
    <h1>Send Commands to Agent</h1>

    <!-- Input for command with 'Enter' key handling -->
    <input 
      bind:value={myMsg} 
      placeholder="Type command here..." 
      on:keydown={handleKeydown} 
    />

    {#if srvMsg}
      <p><strong>Response from Server:</strong> {srvMsg}</p>
    {/if}

    {#if errorMessage}
      <p style="color: red;">{errorMessage}</p>
    {/if}
  </div>
</main>
