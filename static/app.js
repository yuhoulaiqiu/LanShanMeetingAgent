// Global state
let currentMeetingId = null;
let currentSessionId = null;
let currentMeetingContent = null;
let jsonEditor = null;
let summaryJsonEditor = null;
let currentEventSource = null; // 添加全局变量跟踪当前的EventSource

// DOM Elements
const meetingList = document.getElementById('meetingList');
const createMeetingBtn = document.getElementById('createMeetingBtn');
const fileInput = document.getElementById('fileInput');
const noMeetingSelected = document.getElementById('noMeetingSelected');
const meetingDetails = document.getElementById('meetingDetails');
const contentViewer = document.getElementById('contentViewer');
const summaryJsonViewer = document.getElementById('summaryJsonViewer');
const summaryMarkdown = document.getElementById('summaryMarkdown');
const jsonPathInput = document.getElementById('jsonPathInput');
const convertToMarkdownBtn = document.getElementById('convertToMarkdownBtn');
const showJsonBtn = document.getElementById('showJsonBtn');
const chatMessages = document.getElementById('chatMessages');
const chatInput = document.getElementById('chatInput');
const sendMessageBtn = document.getElementById('sendMessageBtn');

// URL State Management
function updateURLState() {
  const params = new URLSearchParams();
  if (currentMeetingId) params.set('meeting', currentMeetingId);
  if (currentTab) params.set('tab', currentTab);

  const newURL = `${window.location.pathname}?${params.toString()}`;
  window.history.pushState({}, '', newURL);
}

function loadURLState() {
  const params = new URLSearchParams(window.location.search);
  const meetingId = params.get('meeting');
  const tab = params.get('tab') || 'content';
  const path = params.get('path');

  if (meetingId) {
    selectMeeting(meetingId);
  }

  if (tab) {
    switchTab(tab);
  }

}

// Tab Management
let currentTab = 'content';
let currentJsonPath = '';

function switchTab(tab) {
  currentTab = tab;
  document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.classList.toggle('active', btn.dataset.tab === tab);
  });
  document.querySelectorAll('.tab-content').forEach(content => {
    content.classList.toggle('active', content.id === `${tab}Tab`);
  });
  updateURLState();
}

// Initialize JSON Editor
function initJsonEditor() {
  const options = {
    mode: 'view',
    modes: ['view', 'code'],
    onModeChange: function (newMode) {
      if (newMode === 'code') {
        jsonEditor.expandAll();
      }
    }
  };
  jsonEditor = new JSONEditor(contentViewer, options);
}

// Initialize Summary JSON Editor
function initSummaryJsonEditor() {
  const options = {
    mode: 'view',
    modes: ['view', 'code'],
    onModeChange: function (newMode) {
      if (newMode === 'code') {
        summaryJsonEditor.expandAll();
      }
    }
  };
  summaryJsonEditor = new JSONEditor(summaryJsonViewer, options);
}

// Get value by JSON path
function getValueByPath(obj, path) {
  const parts = path.split('.');
  let current = obj;

  for (const part of parts) {
    if (part === '$') continue;
    if (current === undefined || current === null) return null;
    current = current[part];
  }

  return current;
}

// Event Listeners
createMeetingBtn.addEventListener('click', () => fileInput.click());
fileInput.addEventListener('change', handleFileUpload);
sendMessageBtn.addEventListener('click', sendMessage);
chatInput.addEventListener('keypress', (e) => {
  if (e.key === 'Enter') sendMessage();
});

convertToMarkdownBtn.addEventListener('click', convertToMarkdown);
showJsonBtn.addEventListener('click', showJson);

// Tab switching
document.querySelectorAll('.tab-btn').forEach(btn => {
  btn.addEventListener('click', () => {
    switchTab(btn.dataset.tab);
  });
});

// Functions
function convertToMarkdown() {
  const path = jsonPathInput.value.trim();
  if (!path) return;

  try {
    const summaryData = summaryJsonEditor.get();
    const value = getValueByPath(summaryData, path);

    if (value === undefined || value === null) {
      alert('No value found at the specified path');
      return;
    }

    // Show raw content
    const content = typeof value === 'string' ? value : JSON.stringify(value, null, 2);

    // Show markdown
    summaryJsonViewer.classList.add('hidden');
    summaryMarkdown.classList.remove('hidden');
    summaryMarkdown.textContent = content;

    // Update URL state
    currentJsonPath = path;
    updateURLState();
  } catch (error) {
    console.error('Error:', error);
    alert('Error converting to markdown');
  }
}

function showJson() {
  summaryJsonViewer.classList.remove('hidden');
  summaryMarkdown.classList.add('hidden');
  currentJsonPath = '';
  updateURLState();
}

async function handleFileUpload(e) {
  const file = e.target.files[0];
  if (!file) return;

  try {
    const content = await file.text();
    const response = await fetch('/meeting', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: content
    });

    if (!response.ok) throw new Error('Failed to create meeting');

    const data = await response.json();
    loadMeetings();
    selectMeeting(data.id);
  } catch (error) {
    console.error('Error:', error);
    alert('Failed to create meeting');
  }
}

async function loadMeetings() {
  try {
    const response = await fetch('/meeting');
    const data = await response.json();

    meetingList.innerHTML = data.meetings.map(meeting => `
            <div class="meeting-item" data-id="${meeting.id}">
                <div class="font-medium">${meeting.content.title || 'Untitled Meeting'}</div>
                <div class="text-sm text-gray-500">${new Date().toLocaleDateString()}</div>
            </div>
        `).join('');

    // Add click handlers to meeting items
    document.querySelectorAll('.meeting-item').forEach(item => {
      item.addEventListener('click', () => selectMeeting(item.dataset.id));
    });
  } catch (error) {
    console.error('Error:', error);
  }
}

async function selectMeeting(meetingId) {
  currentMeetingId = meetingId;
  currentSessionId = `session_${Date.now()}`;

  // Update UI
  document.querySelectorAll('.meeting-item').forEach(item => {
    item.classList.toggle('active', item.dataset.id === meetingId);
  });

  noMeetingSelected.classList.add('hidden');
  meetingDetails.classList.remove('hidden');

  // Load meeting content
  try {
    const response = await fetch('/meeting');
    const data = await response.json();
    const meeting = data.meetings.find(m => m.id === meetingId);
    if (meeting) {
      currentMeetingContent = meeting.content;
      // Update JSON editor
      if (!jsonEditor) {
        initJsonEditor();
      }
      jsonEditor.set(meeting.content);
      jsonEditor.expandAll();
    }
  } catch (error) {
    console.error('Error:', error);
  }

  // Load summary
  try {
    const response = await fetch(`/summary?meeting_id=${meetingId}`);
    const data = await response.json();

    // Initialize summary JSON editor if not exists
    if (!summaryJsonEditor) {
      initSummaryJsonEditor();
    }

    // Update summary JSON editor
    summaryJsonEditor.set(data);
    summaryJsonEditor.expandAll();

    // Reset markdown view
    summaryJsonViewer.classList.remove('hidden');
    summaryMarkdown.classList.add('hidden');
    jsonPathInput.value = '';

    // Update URL state
    updateURLState();
  } catch (error) {
    console.error('Error:', error);
  }

  // Clear chat
  chatMessages.innerHTML = '';
}

async function sendMessage() {
  const message = chatInput.value.trim();
  if (!message || !currentMeetingId || !currentSessionId) return;

  // 如果存在之前的连接，先关闭它
  if (currentEventSource) {
    currentEventSource.close();
    currentEventSource = null;
  }

  // Add user message to chat
  const userMsgID = Math.random().toString(36).substring(2, 15);
  addMessageToChat(userMsgID, message, 'user');
  chatInput.value = '';

  // Start SSE connection and send message
  const eventSource = new EventSource(`/chat?meeting_id=${currentMeetingId}&session_id=${currentSessionId}&message=${encodeURIComponent(message)}`);
  currentEventSource = eventSource; // 保存当前连接的引用
  const assistantMsgID = Math.random().toString(36).substring(2, 15);
  let messageReceived = false; // 标记是否接收到至少一条消息

  eventSource.addEventListener('message', (event) => {
    try {
      messageReceived = true;
      let content = event.data;
      try {
        const jsonData = JSON.parse(event.data);
        if (jsonData.data) {
          content = jsonData.data;
        } else if (jsonData.error) {
          content = `错误: ${jsonData.error}`;
        }
      } catch (e) {
        console.log("使用原始文本内容");
      }
      addMessageToChat(assistantMsgID, content, 'assistant');
    } catch (error) {
      console.error("处理消息时出错:", error);
    }
  });

  eventSource.addEventListener('stop', (event) => {
    console.log("收到停止事件:", event);
    eventSource.close();
    currentEventSource = null; // 清除引用
    console.log("SSE 连接已结束");
  });

  eventSource.onerror = (error) => {
    console.error("SSE连接错误:", error);
    // 检查连接状态：0=连接中, 1=打开, 2=关闭
    if (eventSource.readyState === 2) {
      console.log("连接已关闭");
      eventSource.close();
      currentEventSource = null; // 清除引用
      
      // 如果连接关闭但没有收到任何消息，显示错误提示
      if (!messageReceived) {
        addMessageToChat(assistantMsgID, "与服务器的连接已断开", 'assistant');
      }
    }
  };
}

let msgs = {};

function addMessageToChat(msgID, message, type) {
  if (msgs[msgID]) {
    msgs[msgID].textContent += message;
  } else {
    const messageDiv = document.createElement('div');
    messageDiv.className = `chat-message ${type}`;
    messageDiv.textContent = message;
    chatMessages.appendChild(messageDiv);
    msgs[msgID] = messageDiv;
  }

  chatMessages.scrollTop = chatMessages.scrollHeight;
}

// Initialize
loadMeetings();
loadURLState();
