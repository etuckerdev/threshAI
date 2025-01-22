import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import ChatTerminal from './components/ChatTerminal';

// Import the CSS file
import './components/ChatTerminal.css';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <ChatTerminal />
  </React.StrictMode>
);