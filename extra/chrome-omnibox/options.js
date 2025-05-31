const options = {};
const optionsForm = document.getElementById("options");

optionsForm.baseURL.addEventListener("change", (event) => {
  options.baseURL = event.target.value;
  chrome.storage.sync.set({ options });
});

const data = await chrome.storage.sync.get("options");
console.log(data);
Object.assign(options, data.options);
optionsForm.baseURL.value = options.baseURL;
