let baseURL = "";

chrome.storage.onChanged.addListener((changes, area) => {
  if (area === "sync" && changes.options?.newValue) {
    baseURL = String(changes.options.newValue.baseURL);
  }
  console.log(baseURL);
});

// chrome.action.onClicked.addListener(() => {
//   chrome.tabs.create({ url: "options.html" });
// });

chrome.omnibox.onInputStarted.addListener(function () {
  console.log("💬 onInputStarted");

  chrome.omnibox.setDefaultSuggestion({
    description:
      "Here is a default <match>suggestion</match>. <url>It's <match>url</match> here</url>",
  });
});

chrome.omnibox.onInputChanged.addListener(function (text, suggest) {
  console.log(baseURL);
  console.log("✏️ onInputChanged: " + text);
  suggest([
    {
      content: "https://google.com?biba",
      description: "<url>https://google.com</url>",
      deletable: true,
    },
    {
      content: text + " number two",
      description: "the second entry",
      deletable: true,
    },
  ]);
});

chrome.omnibox.onInputEntered.addListener(function (text, disposition) {
  console.log(
    `✔️ onInputEntered: text -> ${text} | disposition -> ${disposition}`,
  );
});

chrome.omnibox.onInputCancelled.addListener(function () {
  console.log("❌ onInputCancelled");
});

chrome.omnibox.onDeleteSuggestion.addListener(function (text) {
  console.log("⛔ onDeleteSuggestion: " + text);
});

chrome.omnibox.onInputStarted.addListener(function () {
  console.log("started");
});
