const form = document.getElementById('text-form');
const input = document.getElementById('input-text');
const outputSection = document.getElementById('output-section');
const outputText = document.getElementById('output-text');
const diffDiv = document.getElementById('diff');
const historyList = document.getElementById('history-list');
const themeBtn = document.getElementById('theme-toggle-btn');
const themeIcon = document.getElementById('theme-icon');

// История
let history = JSON.parse(localStorage.getItem('history') || '[]');
renderHistory();

form.addEventListener('submit', handleSubmit);
input.addEventListener('keydown', function(e) {
    if (e.key === 'Enter' && !e.shiftKey && !e.ctrlKey && !e.altKey) {
        e.preventDefault();
        form.dispatchEvent(new Event('submit'));
    }
});

// Инициализация темы
(function() {
    const saved = localStorage.getItem('theme');
    if (saved === 'dark') {
        setTheme('dark');
    } else {
        setTheme('light');
    }
})();

themeBtn.addEventListener('click', function() {
    const isDark = document.body.classList.contains('dark');
    setTheme(isDark ? 'light' : 'dark');
});

function setTheme(theme) {
    document.body.classList.toggle('dark', theme === 'dark');
    localStorage.setItem('theme', theme);
    themeIcon.textContent = theme === 'dark' ? '🌙' : '🌞';
}

async function handleSubmit(e) {
    e.preventDefault();
    const inputText = input.value;
    if (!inputText.trim()) return;
    const response = await fetch('/api/process', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text: inputText })
    });
    const data = await response.json();
    outputSection.style.display = 'block';
    outputText.textContent = data.result;
    diffDiv.innerHTML = diffStrings(inputText, data.result);
    // Добавить в историю
    const entry = {
        input: inputText,
        output: data.result,
        diff: diffStrings(inputText, data.result)
    };
    history.unshift(entry);
    if (history.length > 10) history = history.slice(0, 10);
    localStorage.setItem('history', JSON.stringify(history));
    renderHistory();
}

function renderHistory() {
    historyList.innerHTML = '';
    if (!history.length) {
        historyList.innerHTML = '<div>История пуста</div>';
        return;
    }
    history.forEach(entry => {
        const div = document.createElement('div');
        div.innerHTML = `<div class="history-input"><b>Ввод:</b> ${escapeHtml(entry.input)}</div><div class="history-output"><b>Вывод:</b> ${escapeHtml(entry.output)}</div><div class="history-diff">${entry.diff}</div>`;
        historyList.appendChild(div);
    });
}

// Простая функция сравнения строк (построчно, с подсветкой изменений)
function diffStrings(a, b) {
    const aLines = a.split(/\r?\n/);
    const bLines = b.split(/\r?\n/);
    let result = '';
    const maxLen = Math.max(aLines.length, bLines.length);
    for (let i = 0; i < maxLen; i++) {
        const orig = aLines[i] || '';
        const corr = bLines[i] || '';
        if (orig === corr) {
            result += `<div>${escapeHtml(orig)}</div>`;
        } else {
            result += `<div><span class="diff-removed">${escapeHtml(orig)}</span> &rarr; <span class="diff-added">${escapeHtml(corr)}</span></div>`;
        }
    }
    return result;
}
function escapeHtml(text) {
    return text.replace(/[&<>"']/g, function(m) {
        return ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'})[m];
    });
} 