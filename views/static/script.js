// Keep DEK in JS memory (and sessionStorage for page refresh in same tab)
window.passman = window.passman || {};
try {
    window.passman.dek = window.passman.dek || sessionStorage.getItem('dek');
} catch (_) {
    window.passman.dek = window.passman.dek || null;
}

if (window.htmx) {
    document.addEventListener('htmx:configRequest', function (event) {
        if (window.passman && window.passman.dek) {
            event.detail.headers['X-DEK'] = window.passman.dek;
        }
    });
}

async function copyPass(pass) {
    const text = document.getElementById(pass).innerText;
    const type = "text/plain";
    const clipboardItemData = { [type]: new Blob([text], { type }) };
    const clipboardItem = new ClipboardItem(clipboardItemData);
    await navigator.clipboard.write([clipboardItem]);
}

document.addEventListener('DOMContentLoaded', function() {
    const addBtn = document.getElementById('addNewPasswordBtn');
    const cancelBtn = document.getElementById('cancelBtn');
    const formContainer = document.getElementById('newCredentialForm');
    const form = formContainer?.querySelector('form');

    console.log('Elements found:', { addBtn, cancelBtn, formContainer, form });
    
    if (!addBtn || !cancelBtn || !formContainer || !form) {
        console.warn('Missing form elements');
        return;
    }

    addBtn.addEventListener('click', function() {
        console.log
        formContainer.classList.remove('hidden');
        addBtn.classList.add('hidden');
    });

    cancelBtn.addEventListener('click', function() {
        formContainer.classList.add('hidden');
        addBtn.classList.remove('hidden');
        form.reset();
    });

    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const existingWebsite = document.getElementById('existingWebsite').value;
        const newWebsite = document.getElementById('newWebsite').value.trim();
        const website = newWebsite !== '' ? newWebsite : existingWebsite;
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
        const headers = {
                'Content-Type': 'application/json',
            };
        if (window.passman && window.passman.dek) {
            headers['X-DEK'] = window.passman.dek;
        }

        const response = await fetch('/api/newcredential', {
                method: 'POST',
                headers,
                body: JSON.stringify({ website, username, password })
            });

            if (response.ok) {
                formContainer.classList.add('hidden');
                addBtn.classList.remove('hidden');
                form.reset();
                window.location.reload();
            } else {
                const errorData = await response.json().catch(() => ({ message: 'Failed to save credential' }));
                alert(errorData.message || 'Failed to save credential');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('Failed to save credential');
        }
    });
});
