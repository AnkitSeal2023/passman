async function copyPass(pass) {
    const text = document.getElementById(pass).innerText;
    const type = "text/plain";
    const clipboardItemData = { [type]: new Blob([text], { type }) };
    const clipboardItem = new ClipboardItem(clipboardItemData);
    await navigator.clipboard.write([clipboardItem]);
}