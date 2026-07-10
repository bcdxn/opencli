(function (global) {
  "use strict";

  if (!global || !global.document) {
    return;
  }

  var EMBED_MARKUP = "\n\n\u003cstyle\u003e\n \n\n \n.ocli[data-theme=\"light\"] {\n  --ocli-bg:           #ffffff;\n  --ocli-bg-sidebar:   #f7f9fc;\n  --ocli-bg-code:      #eef1f6;\n  --ocli-bg-hover:     #eef2f8;\n  --ocli-border:       #dde3ec;\n  --ocli-border-light: #edf0f5;\n  --ocli-text:         #3a4a5c;\n  --ocli-text-muted:   #7a8aa0;\n  --ocli-text-active:  #0d1f2d;\n  --ocli-accent:       #5469d4;\n  --ocli-accent-soft:  #eef0fb;\n  --ocli-badge-bg:     #e8edf8;\n  --ocli-badge-text:   #3d5caa;\n  --ocli-required-bg:  #fef2f2;\n  --ocli-required-text:#c0392b;\n  --ocli-group-bg:     #f0f4f8;\n  --ocli-group-text:   #5a6a7e;\n  --ocli-sticky-shadow:0 2px 8px rgba(0,0,0,0.06);\n}\n\n.ocli[data-theme=\"dark\"] {\n  --ocli-bg:           #2e3440;\n  --ocli-bg-sidebar:   #262c38;\n  --ocli-bg-code:      #3b4252;\n  --ocli-bg-hover:     #353d4f;\n  --ocli-border:       #4c566a;\n  --ocli-border-light: #3e4759;\n  --ocli-text:         #c8d3e6;\n  --ocli-text-muted:   #7b8fa8;\n  --ocli-text-active:  #eceff4;\n  --ocli-accent:       #81a1c1;\n  --ocli-accent-soft:  #2d3b4f;\n  --ocli-badge-bg:     #2c3d55;\n  --ocli-badge-text:   #88c0d0;\n  --ocli-required-bg:  #3a2020;\n  --ocli-required-text:#f08080;\n  --ocli-group-bg:     #3b4252;\n  --ocli-group-text:   #8a9ab5;\n  --ocli-sticky-shadow:0 2px 8px rgba(0,0,0,0.25);\n}\n\n \n.ocli, .ocli *, .ocli *::before, .ocli *::after {\n  box-sizing: border-box;\n}\n\n.ocli {\n  background: var(--ocli-bg);\n  color: var(--ocli-text);\n  display: flex;\n  flex-direction: row;\n  font-family: -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, Oxygen,\n               Ubuntu, Cantarell, \"Helvetica Neue\", sans-serif;\n  font-size: 15px;\n  line-height: 1.6;\n  height: 100%;\n  min-width: 0;\n}\n\n \n.ocli__sidebar {\n  background: var(--ocli-bg-sidebar);\n  border-right: 1px solid var(--ocli-border);\n  display: flex;\n  flex-direction: column;\n  flex-shrink: 0;\n  height: 100vh;\n  overflow-y: auto;\n  position: sticky;\n  top: 0;\n  width: 260px;\n}\n\n.ocli__sidebar::-webkit-scrollbar { width: 4px; }\n.ocli__sidebar::-webkit-scrollbar-thumb { background: var(--ocli-border); border-radius: 2px; }\n\n.ocli__sidebar-header {\n  align-items: stretch;\n  border-bottom: 1px solid var(--ocli-border);\n  display: flex;\n  flex-direction: column;\n  gap: 6px;\n  justify-content: center;\n  min-height: 60px;\n  padding: 12px 16px;\n}\n\n.ocli__sidebar-header-top {\n  align-items: center;\n  display: flex;\n  gap: 10px;\n  justify-content: space-between;\n}\n\n.ocli__sidebar-title {\n  color: var(--ocli-text-active);\n  font-size: 15px;\n  font-weight: 700;\n  line-height: 1.3;\n}\n\n.ocli__sidebar-version {\n  color: var(--ocli-text-muted);\n  font-size: 11px;\n  font-weight: 500;\n}\n\n.ocli-theme-toggle {\n  align-items: center;\n  background: var(--ocli-bg-hover);\n  border: 1px solid var(--ocli-border);\n  border-radius: 999px;\n  color: var(--ocli-text-muted);\n  cursor: pointer;\n  display: inline-flex;\n  flex-shrink: 0;\n  height: 28px;\n  justify-content: flex-start;\n  padding: 3px;\n  transition: background 0.2s ease, border-color 0.2s ease;\n  width: 52px;\n}\n\n.ocli-theme-toggle:hover {\n  border-color: var(--ocli-accent);\n}\n\n.ocli-theme-toggle:focus-visible {\n  outline: 2px solid var(--ocli-accent);\n  outline-offset: 2px;\n}\n\n.ocli-theme-toggle__thumb {\n  align-items: center;\n  background: var(--ocli-bg);\n  border-radius: 999px;\n  box-shadow: 0 1px 3px rgba(0,0,0,0.25);\n  display: inline-flex;\n  height: 20px;\n  justify-content: center;\n  position: relative;\n  transform: translateX(0);\n  transition: transform 0.2s ease, background 0.2s ease;\n  width: 20px;\n}\n\n.ocli-theme-toggle__icon {\n  fill: none;\n  height: 12px;\n  left: 50%;\n  position: absolute;\n  stroke: currentColor;\n  stroke-linecap: round;\n  stroke-linejoin: round;\n  stroke-width: 1.8;\n  top: 50%;\n  transform: translate(-50%, -50%);\n  transition: opacity 0.15s ease;\n  width: 12px;\n}\n\n.ocli-theme-toggle__icon--sun {\n  opacity: 1;\n}\n\n.ocli-theme-toggle__icon--moon {\n  opacity: 0;\n}\n\n.ocli[data-theme=\"light\"] .ocli-theme-toggle {\n  background: #d9dee7;\n}\n\n.ocli[data-theme=\"light\"] .ocli-theme-toggle__thumb {\n  background: #ffffff;\n  color: #d59a22;\n}\n\n.ocli[data-theme=\"dark\"] .ocli-theme-toggle {\n  background: #5d6577;\n}\n\n.ocli[data-theme=\"dark\"] .ocli-theme-toggle__thumb {\n  background: #1f2531;\n  color: #d6dded;\n  transform: translateX(24px);\n}\n\n.ocli[data-theme=\"dark\"] .ocli-theme-toggle__icon--sun {\n  opacity: 0;\n}\n\n.ocli[data-theme=\"dark\"] .ocli-theme-toggle__icon--moon {\n  opacity: 1;\n}\n\n.ocli__sidebar-body {\n  flex: 1;\n  overflow-y: auto;\n  padding: 12px 0 24px;\n}\n\n.ocli__nav-section {\n  margin: 16px 0 4px;\n  padding: 0 16px;\n}\n\n.ocli__nav-section-label {\n  color: var(--ocli-text-muted);\n  font-size: 10px;\n  font-weight: 700;\n  letter-spacing: 0.08em;\n  padding: 4px 0;\n  text-transform: uppercase;\n}\n\n.ocli__nav-list {\n  list-style: none;\n  margin: 0;\n  padding: 0;\n}\n\n.ocli__nav-list--nested {\n  border-left: 2px solid var(--ocli-border-light);\n  margin-left: 10px;\n  padding-left: 8px;\n}\n\n.ocli__nav-item {\n  margin: 0;\n  padding: 0;\n}\n\n.ocli__nav-link {\n  align-items: center;\n  border-radius: 5px;\n  color: var(--ocli-text);\n  cursor: pointer;\n  display: flex;\n  font-size: 13px;\n  gap: 4px;\n  justify-content: space-between;\n  padding: 5px 8px;\n  text-decoration: none;\n  transition: background 0.1s, color 0.1s;\n}\n\n.ocli__nav-link:hover {\n  background: var(--ocli-bg-hover);\n  color: var(--ocli-text-active);\n}\n\n.ocli__nav-link.is-active {\n  background: var(--ocli-accent-soft);\n  color: var(--ocli-accent);\n  font-weight: 600;\n}\n\n.ocli__nav-link .ocli__nav-toggle {\n  display: inline-block;\n  font-size: 10px;\n  opacity: 0.6;\n  transition: transform 0.15s;\n}\n\n.ocli__nav-item.is-open \u003e .ocli__nav-link .ocli__nav-toggle {\n  transform: rotate(90deg);\n}\n\n.ocli__nav-sublist {\n  display: none;\n  list-style: none;\n  margin: 2px 0 4px 12px;\n  padding: 0;\n}\n\n.ocli__nav-item.is-open \u003e .ocli__nav-sublist {\n  display: block;\n}\n\n \n.ocli__main {\n  flex: 1;\n  height: 100%;\n  min-width: 0;\n  overflow-y: auto;\n}\n\n.ocli__main::-webkit-scrollbar { width: 6px; }\n.ocli__main::-webkit-scrollbar-thumb { background: var(--ocli-border); border-radius: 3px; }\n\n \n.ocli-section {\n  border-bottom: 1px solid var(--ocli-border);\n}\n\n.ocli-section__sticky {\n  align-items: center;\n  background: var(--ocli-bg);\n  border-bottom: 1px solid transparent;\n  display: flex;\n  flex-wrap: wrap;\n  gap: 8px;\n  min-height: 60px;\n  padding: 0 24px;\n  position: sticky;\n  top: -1px;\n  transition: border-color 0.15s, box-shadow 0.15s;\n  z-index: 10;\n}\n\n.ocli-section__sticky.is-stuck {\n  border-bottom-color: var(--ocli-border);\n  box-shadow: var(--ocli-sticky-shadow);\n}\n\n.ocli-section__sticky-title {\n  color: var(--ocli-text-active);\n  font-size: 16px;\n  font-weight: 600;\n}\n\n.ocli-section__cmd {\n  background: var(--ocli-bg-code);\n  border-radius: 5px;\n  color: var(--ocli-text-active);\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  font-size: 14px;\n  font-weight: 600;\n  padding: 3px 8px;\n  margin-left: 6px;\n}\n\n.ocli-section__body {\n  padding: 0px 24px 28px 30px;\n}\n\n.ocli-section__intro-grid {\n  display: grid;\n  gap: 24px;\n  grid-template-columns: 1fr 1fr;\n  padding: 20px 24px 28px;\n}\n\n@media (max-width: 720px) {\n  .ocli-section__intro-grid {\n    grid-template-columns: 1fr;\n  }\n}\n\n \n.ocli-badges {\n  align-items: center;\n  display: flex;\n  flex-wrap: wrap;\n  gap: 5px;\n}\n\n.ocli-badge {\n  border-radius: 4px;\n  font-size: 11px;\n  font-weight: 600;\n  letter-spacing: 0.02em;\n  padding: 2px 7px;\n}\n\n.ocli-badge--default {\n  background: var(--ocli-badge-bg);\n  color: var(--ocli-badge-text);\n}\n\n.ocli-badge--group {\n  background: var(--ocli-group-bg);\n  color: var(--ocli-group-text);\n}\n\n.ocli-badge--required {\n  background: var(--ocli-required-bg);\n  color: var(--ocli-required-text);\n}\n\n \n.ocli-section__body h2, .ocli-section__intro-grid h2 {\n  color: var(--ocli-text-active);\n  font-size: 14px;\n  font-weight: 700;\n  margin: 0 0 8px;\n  text-transform: uppercase;\n  letter-spacing: 0.04em;\n}\n\n.ocli-summary {\n  color: var(--ocli-text-active);\n  font-size: 15px;\n  font-weight: 500;\n  margin-bottom: 12px;\n}\n\n \n.ocli-md p {\n  color: var(--ocli-text);\n  line-height: 1.7;\n  margin: 0 0 12px;\n}\n.ocli-md p:last-child { margin-bottom: 0; }\n.ocli-md a { color: var(--ocli-accent); text-decoration: none; }\n.ocli-md a:hover { text-decoration: underline; }\n.ocli-md strong { color: var(--ocli-text-active); font-weight: 600; }\n.ocli-md em { font-style: italic; }\n.ocli-md ul, .ocli-md ol { margin: 8px 0 12px 20px; }\n.ocli-md li { margin: 4px 0; }\n.ocli-md code {\n  background: var(--ocli-bg-code);\n  border-radius: 4px;\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  font-size: 0.875em;\n  padding: 2px 5px;\n}\n.ocli-md pre {\n  background: var(--ocli-bg-code);\n  border-radius: 6px;\n  margin: 12px 0;\n  overflow-x: auto;\n  padding: 14px 16px;\n}\n.ocli-md pre code {\n  background: none;\n  font-size: 13px;\n  padding: 0;\n}\n.ocli-md h1, .ocli-md h2, .ocli-md h3 {\n  color: var(--ocli-text-active);\n  font-weight: 600;\n  margin: 16px 0 8px;\n}\n\n \n.ocli-meta-box {\n  margin-bottom: 16px;\n}\n\n.ocli-meta-box__label {\n  color: var(--ocli-text-muted);\n  font-size: 11px;\n  font-weight: 700;\n  letter-spacing: 0.07em;\n  margin-bottom: 6px;\n  text-transform: uppercase;\n}\n\n.ocli-meta-box__value {\n  color: var(--ocli-text);\n  font-size: 14px;\n}\n\n.ocli-meta-box__value a {\n  color: var(--ocli-accent);\n  text-decoration: none;\n}\n\n.ocli-meta-box__value a:hover { text-decoration: underline; }\n\n \n.ocli-tabs__tablist {\n  border-bottom: 1px solid var(--ocli-border);\n  display: flex;\n  flex-wrap: wrap;\n  gap: 0;\n  margin-bottom: 20px;\n}\n\n.ocli-tabs__tab {\n  background: none;\n  border: none;\n  border-bottom: 2px solid transparent;\n  color: var(--ocli-text-muted);\n  cursor: pointer;\n  font-family: inherit;\n  font-size: 14px;\n  font-weight: 500;\n  margin-bottom: -1px;\n  padding: 10px 16px;\n  transition: color 0.15s, border-color 0.15s;\n}\n\n.ocli-tabs__tab:hover {\n  color: var(--ocli-text-active);\n}\n\n.ocli-tabs__tab.is-active {\n  border-bottom-color: var(--ocli-accent);\n  color: var(--ocli-accent);\n}\n\n.ocli-tabs__panel {\n  display: none;\n}\n\n.ocli-tabs__panel.is-active {\n  display: block;\n}\n\n.ocli-tabs__panel-desc {\n  color: var(--ocli-text);\n  font-size: 14px;\n  margin: 0 0 12px;\n}\n\n.ocli-tabs__panel-url {\n  font-size: 14px;\n  margin-top: 10px;\n}\n\n.ocli-tabs__panel-url a {\n  color: var(--ocli-accent);\n}\n\n \n.ocli-code-block {\n  background: var(--ocli-bg-code);\n  border-radius: 6px;\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  font-size: 13px;\n  line-height: 1.6;\n  overflow-x: auto;\n  padding: 12px 14px;\n}\n\n \n.ocli-code-block--terminal {\n  background: none;\n  border: 1px solid var(--ocli-border);\n  border-radius: 8px;\n  overflow: hidden;\n  padding: 0;\n}\n\n.ocli-code-block__header {\n  align-items: center;\n  background: #1a2133;\n  border-bottom: 1px solid rgba(255,255,255,0.08);\n  display: flex;\n  padding: 8px 14px;\n}\n\n.ocli-code-block__label {\n  color: rgba(255,255,255,0.55);\n  font-size: 12px;\n  font-weight: 500;\n  letter-spacing: 0.02em;\n}\n\n.ocli-code-block__body {\n  background: #1e2d45;\n  color: #a8d8f0;\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  font-size: 13px;\n  line-height: 1.6;\n  margin: 0;\n  overflow-x: auto;\n  padding: 16px 18px;\n  white-space: pre;\n}\n\n \n.ocli-table {\n  border-collapse: collapse;\n  font-size: 14px;\n  width: 100%;\n}\n\n.ocli-table th {\n  background: var(--ocli-bg-sidebar);\n  border: 1px solid var(--ocli-border);\n  color: var(--ocli-text-muted);\n  font-size: 11px;\n  font-weight: 700;\n  letter-spacing: 0.05em;\n  padding: 8px 12px;\n  text-align: left;\n  text-transform: uppercase;\n}\n\n.ocli-table td {\n  border: 1px solid var(--ocli-border);\n  padding: 8px 12px;\n  color: var(--ocli-text);\n}\n\n.ocli-table code {\n  background: var(--ocli-bg-code);\n  border-radius: 3px;\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  font-size: 0.875em;\n  padding: 1px 4px;\n}\n\n \n.ocli-exit-code {\n  align-items: flex-start;\n  border-bottom: 1px solid var(--ocli-border-light);\n  display: flex;\n  gap: 12px;\n  padding: 12px 0;\n}\n\n.ocli-exit-code:last-child { border-bottom: none; }\n\n.ocli-exit-code__code {\n  background: var(--ocli-bg-code);\n  border-radius: 4px;\n  color: var(--ocli-text-active);\n  flex-shrink: 0;\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  font-size: 13px;\n  font-weight: 600;\n  min-width: 28px;\n  padding: 2px 7px;\n  text-align: center;\n}\n\n.ocli-exit-code__details {}\n\n.ocli-exit-code__status {\n  background: var(--ocli-badge-bg);\n  border-radius: 4px;\n  color: var(--ocli-badge-text);\n  display: inline-block;\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  font-size: 11px;\n  font-weight: 600;\n  margin-bottom: 4px;\n  padding: 2px 6px;\n}\n\n.ocli-exit-code__summary {\n  color: var(--ocli-text);\n  font-size: 14px;\n}\n\n.ocli-exit-code__desc {\n  color: var(--ocli-text-muted);\n  font-size: 13px;\n  margin-top: 4px;\n}\n\n \n.ocli-params {\n  margin-top: 24px;\n}\n\n.ocli-params__heading {\n  border-bottom: 1px solid var(--ocli-border);\n  color: var(--ocli-text-muted);\n  font-size: 11px;\n  font-weight: 700;\n  letter-spacing: 0.07em;\n  margin-bottom: 4px;\n  padding-bottom: 8px;\n  text-transform: uppercase;\n}\n\n.ocli-param {\n  border-bottom: 1px solid var(--ocli-border-light);\n  padding: 14px 0;\n  scroll-margin-top: 60px;\n}\n\n.ocli-param:last-child { border-bottom: none; }\n\n.ocli-param__header {\n  align-items: center;\n  display: flex;\n  flex-wrap: wrap;\n  gap: 8px;\n  margin-bottom: 6px;\n}\n\n.ocli-param__name {\n  color: var(--ocli-text-active);\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  font-size: 13px;\n  font-weight: 600;\n  padding: 2px 7px 2px 4px;\n}\n\n \n.ocli-param__anchor,\n.ocli-section__cmd-anchor {\n  align-items: center;\n  color: inherit;\n  display: inline-flex;\n  gap: 6px;\n  position: relative;\n}\n\n.ocli-deeplink-btn {\n  background: none;\n  border: none;\n  cursor: pointer;\n  left: -20px;\n  padding: 0 6px;\n  position: absolute;\n  top: 50%;\n  transform: translateY(-50%);\n}\n\n.ocli-deeplink-icon {\n  color: var(--ocli-accent);\n  display: block;\n  opacity: 0;\n  transition: opacity 0.15s;\n}\n\n.ocli-param__anchor:hover .ocli-deeplink-icon,\n.ocli-section__cmd-anchor:hover .ocli-deeplink-icon {\n  opacity: 1;\n}\n\n.ocli-param__summary {\n  color: var(--ocli-text);\n  font-size: 14px;\n  margin-bottom: 4px;\n}\n\n.ocli-param__desc {\n  color: var(--ocli-text);\n  font-size: 14px;\n  margin-top: 6px;\n}\n\n.ocli-param__desc p { margin: 0 0 4px; }\n\n \n.ocli-choices {\n  margin-top: 10px;\n}\n\n.ocli-choices__label {\n  color: var(--ocli-text-muted);\n  font-size: 12px;\n  font-weight: 600;\n  margin-bottom: 6px;\n  text-transform: uppercase;\n  letter-spacing: 0.05em;\n}\n\n.ocli-choices__list {\n  list-style: none;\n  margin: 0;\n  padding: 0;\n}\n\n.ocli-choices__item {\n  align-items: baseline;\n  display: flex;\n  font-size: 13px;\n  gap: 8px;\n  padding: 3px 0;\n}\n\n.ocli-choices__item code {\n  background: var(--ocli-bg-code);\n  border-radius: 3px;\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  flex-shrink: 0;\n  padding: 1px 5px;\n}\n\n.ocli-choices__item-desc { color: var(--ocli-text-muted); }\n\n \n.ocli-alt-sources {\n  margin-top: 10px;\n}\n\n.ocli-alt-sources__label {\n  color: var(--ocli-text-muted);\n  font-size: 12px;\n  font-weight: 600;\n  margin-bottom: 6px;\n  text-transform: uppercase;\n  letter-spacing: 0.05em;\n}\n\n.ocli-alt-sources__list {\n  list-style: none;\n  margin: 0;\n  padding: 0;\n}\n\n.ocli-alt-sources__item {\n  font-size: 13px;\n  padding: 2px 0;\n}\n\n.ocli-alt-sources__item code {\n  background: var(--ocli-bg-code);\n  border-radius: 3px;\n  font-family: \"SFMono-Regular\", Consolas, \"Liberation Mono\", Menlo, monospace;\n  padding: 1px 5px;\n}\n\n \n.ocli-footer {\n  color: var(--ocli-text-muted);\n  font-size: 12px;\n  padding: 24px;\n  text-align: center;\n}\n\n.ocli-footer a { color: var(--ocli-accent); text-decoration: none; }\n.ocli-footer a:hover { text-decoration: underline; }\n\n \n@media (max-width: 900px) {\n  .ocli__sidebar {\n    display: none;\n  }\n}\n\u003c/style\u003e\n\n\n\u003cdiv class=\"ocli\" data-theme=\"light\"\u003e\n  \n\u003cnav class=\"ocli__sidebar\" aria-label=\"Navigation\"\u003e\n  \u003cdiv class=\"ocli__sidebar-header\"\u003e\n    \u003cdiv class=\"ocli__sidebar-header-top\"\u003e\n      \u003cspan class=\"ocli__sidebar-title\"\u003eOpenCLI\u003c/span\u003e\n      \u003cbutton\n        class=\"ocli-theme-toggle\"\n        type=\"button\"\n        data-theme-toggle\n        aria-label=\"Switch to dark theme\"\n        aria-pressed=\"false\"\n      \u003e\n        \u003cspan class=\"ocli-theme-toggle__thumb\" aria-hidden=\"true\"\u003e\n          \u003csvg class=\"ocli-theme-toggle__icon ocli-theme-toggle__icon--sun\" viewBox=\"0 0 24 24\" focusable=\"false\"\u003e\n            \u003ccircle cx=\"12\" cy=\"12\" r=\"4\"\u003e\u003c/circle\u003e\n            \u003cpath d=\"M12 2v2.2M12 19.8V22M4.9 4.9l1.6 1.6M17.5 17.5l1.6 1.6M2 12h2.2M19.8 12H22M4.9 19.1l1.6-1.6M17.5 6.5l1.6-1.6\"\u003e\u003c/path\u003e\n          \u003c/svg\u003e\n          \u003csvg class=\"ocli-theme-toggle__icon ocli-theme-toggle__icon--moon\" viewBox=\"0 0 24 24\" focusable=\"false\"\u003e\n            \u003cpath d=\"M20.4 14.8A8.7 8.7 0 0 1 9.2 3.6a8.6 8.6 0 1 0 11.2 11.2z\"\u003e\u003c/path\u003e\n          \u003c/svg\u003e\n        \u003c/span\u003e\n      \u003c/button\u003e\n    \u003c/div\u003e\n    \u003cspan class=\"ocli__sidebar-version\"\u003ev1.0.0-alpha.12\u003c/span\u003e\n  \u003c/div\u003e\n  \u003cdiv class=\"ocli__sidebar-body\"\u003e\n    \u003cdiv class=\"ocli__nav-section\"\u003e\n      \u003cdiv class=\"ocli__nav-section-label\"\u003eOverview\u003c/div\u003e\n      \u003cul class=\"ocli__nav-list\"\u003e\n        \u003cli class=\"ocli__nav-item\"\u003e\n          \u003ca class=\"ocli__nav-link\" href=\"#intro\" data-section=\"intro\"\u003eIntroduction\u003c/a\u003e\n        \u003c/li\u003e\n        \u003cli class=\"ocli__nav-item\"\u003e\n          \u003ca class=\"ocli__nav-link\" href=\"#install\" data-section=\"install\"\u003eInstallation\u003c/a\u003e\n        \u003c/li\u003e\n          \u003cli class=\"ocli__nav-item\"\u003e\n            \u003ca class=\"ocli__nav-link\" href=\"#config\" data-section=\"config\"\u003eConfiguration\u003c/a\u003e\n          \u003c/li\u003e\n          \u003cli class=\"ocli__nav-item\"\u003e\n            \u003ca class=\"ocli__nav-link\" href=\"#exit-codes\" data-section=\"exit-codes\"\u003eExit Codes\u003c/a\u003e\n          \u003c/li\u003e\n      \u003c/ul\u003e\n    \u003c/div\u003e\n    \u003cdiv class=\"ocli__nav-section\"\u003e\n      \u003cdiv class=\"ocli__nav-section-label\"\u003eCommands\u003c/div\u003e\n      \u003cul class=\"ocli__nav-list\"\u003e\n\u003cli class=\"ocli__nav-item\"\u003e\n  \u003ca class=\"ocli__nav-link\" href=\"#ocli-check\" data-section=\"ocli-check\"\u003e\n    \u003cspan\u003echeck\u003c/span\u003e\n  \u003c/a\u003e\n\u003c/li\u003e\n\n\u003cli class=\"ocli__nav-item ocli__nav-item--group\"\u003e\n  \u003ca class=\"ocli__nav-link\" href=\"#ocli-gen\" data-section=\"ocli-gen\"\u003e\n    \u003cspan\u003egen\u003c/span\u003e\n    \u003cspan class=\"ocli__nav-toggle\"\u003e▸\u003c/span\u003e\n  \u003c/a\u003e\n  \u003cul class=\"ocli__nav-sublist\"\u003e\n\u003cli class=\"ocli__nav-item\"\u003e\n  \u003ca class=\"ocli__nav-link\" href=\"#ocli-gen-docs\" data-section=\"ocli-gen-docs\"\u003e\n    \u003cspan\u003edocs\u003c/span\u003e\n  \u003c/a\u003e\n\u003c/li\u003e\n\n\u003cli class=\"ocli__nav-item\"\u003e\n  \u003ca class=\"ocli__nav-link\" href=\"#ocli-gen-cli\" data-section=\"ocli-gen-cli\"\u003e\n    \u003cspan\u003ecli\u003c/span\u003e\n  \u003c/a\u003e\n\u003c/li\u003e\n\n  \u003c/ul\u003e\n\u003c/li\u003e\n\n      \u003c/ul\u003e\n    \u003c/div\u003e\n  \u003c/div\u003e\n\u003c/nav\u003e\n\n  \u003cdiv class=\"ocli__main\"\u003e\n    \n\u003csection class=\"ocli-section\" id=\"intro\"\u003e\n  \u003cdiv class=\"ocli-section__sticky\"\u003e\n    \u003cspan class=\"ocli-section__sticky-title\"\u003eIntroduction\u003c/span\u003e\n  \u003c/div\u003e\n  \u003cdiv class=\"ocli-section__intro-grid\"\u003e\n    \u003cdiv class=\"ocli-section__intro-left\"\u003e\n      \u003cimg src=\"https://img.shields.io/badge/OpenCLI_Spec-Compliant-brightgreen?link=https%3A%2F%2Fgithub.com%2Fbcdxn%2Fopencli\"/\u003e\n      \u003cp class=\"ocli-summary\"\u003eA CLI for working with OpenCLI Specs\u003c/p\u003e\n      \u003cdiv class=\"ocli-md\"\u003e\u003cp\u003e\u003ccode\u003eocli\u003c/code\u003e is a command line interface designed to make working with \u003ca href=\"https://github.com/bcdxn/opencli/tree/main\"\u003eOpenCLI\nSpec documents\u003c/a\u003e easier. It provides a number of capabilities, including:\u003c/p\u003e\n\u003cul\u003e\n\u003cli\u003evalidating OpenCLI Spec documents\u003c/li\u003e\n\u003cli\u003egenerating Documentation from OpenCLI Spec documents\u003c/li\u003e\n\u003c/ul\u003e\n\u003cp\u003eThe commands are documented below. You can also find out more about each\ncommand using the contextual \u003ccode\u003e--help\u003c/code\u003e flag. e.g.:\u003c/p\u003e\n\u003cpre\u003e\u003ccode class=\"language-sh\"\u003eocli gen --help\n\u003c/code\u003e\u003c/pre\u003e\n\u003c/div\u003e\n    \u003c/div\u003e\n    \u003cdiv class=\"ocli-section__intro-right\"\u003e\n      \u003cdiv class=\"ocli-meta-box\"\u003e\n        \u003cdiv class=\"ocli-meta-box__label\"\u003eLicense\u003c/div\u003e\n        \u003cdiv class=\"ocli-meta-box__value\"\u003e\n          \u003ca href=\"https://spdx.org/licenses/MIT.html\"\u003eMIT\u003c/a\u003e\n        \u003c/div\u003e\n      \u003c/div\u003e\n      \u003cdiv class=\"ocli-meta-box\"\u003e\n        \u003cdiv class=\"ocli-meta-box__label\"\u003eContact\u003c/div\u003e\n        \u003cdiv class=\"ocli-meta-box__value\"\u003e\u003cdiv\u003eOpenCLI Team\u003c/div\u003e\u003cdiv\u003e\u003ca href=\"https://github.com/bcdxn/opencli\"\u003ehttps://github.com/bcdxn/opencli\u003c/a\u003e\u003c/div\u003e\n        \u003c/div\u003e\n      \u003c/div\u003e\n      \u003cdiv class=\"ocli-meta-box\"\u003e\n        \u003cdiv class=\"ocli-meta-box__label\"\u003eBinary\u003c/div\u003e\n        \u003cdiv class=\"ocli-meta-box__value\"\u003e\u003ccode style=\"font-family:monospace;background:var(--ocli-bg-code);padding:2px 6px;border-radius:4px;\"\u003eocli\u003c/code\u003e\u003c/div\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n  \u003c/div\u003e\n\u003c/section\u003e\n\n\u003csection class=\"ocli-section\" id=\"install\"\u003e\n  \u003cdiv class=\"ocli-section__sticky\"\u003e\n    \u003cspan class=\"ocli-section__sticky-title\"\u003eInstallation\u003c/span\u003e\n  \u003c/div\u003e\n  \u003cdiv class=\"ocli-section__body\"\u003e\n    \u003cdiv class=\"ocli-tabs\"\u003e\n      \u003cdiv class=\"ocli-tabs__tablist\" role=\"tablist\"\u003e\n        \u003cbutton class=\"ocli-tabs__tab is-active\"\n                role=\"tab\"\n                aria-selected=\"true\"\n                aria-controls=\"install-panel-0\"\n                data-tab=\"install-panel-0\"\u003eBinary Download\u003c/button\u003e\n        \u003cbutton class=\"ocli-tabs__tab\"\n                role=\"tab\"\n                aria-selected=\"false\"\n                aria-controls=\"install-panel-1\"\n                data-tab=\"install-panel-1\"\u003eGo Install\u003c/button\u003e\n      \u003c/div\u003e\n      \u003cdiv class=\"ocli-tabs__panel is-active\"\n           id=\"install-panel-0\"\n           role=\"tabpanel\"\u003e\n        \u003cp class=\"ocli-tabs__panel-desc\"\u003eYou can visit the GitHub releases page of the project to download pre-compiled binaries\nappropriate for your system.\n\u003c/p\u003e\n        \u003cdiv class=\"ocli-tabs__panel-url\"\u003e\u003ca href=\"https://github.com/bcdxn/opencli/releases\"\u003ehttps://github.com/bcdxn/opencli/releases\u003c/a\u003e\u003c/div\u003e\n      \u003c/div\u003e\n      \u003cdiv class=\"ocli-tabs__panel\"\n           id=\"install-panel-1\"\n           role=\"tabpanel\"\u003e\n        \u003cdiv class=\"ocli-code-block ocli-code-block--terminal\"\u003e\n          \u003cdiv class=\"ocli-code-block__header\"\u003e\n            \u003cspan class=\"ocli-code-block__label\"\u003eCommand Line\u003c/span\u003e\n          \u003c/div\u003e\n          \u003cpre class=\"ocli-code-block__body\"\u003ego install github.com/bcdxn/opencli/cmd/cobra@latest\u003c/pre\u003e\n        \u003c/div\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n  \u003c/div\u003e\n\u003c/section\u003e\n\n      \n\u003csection class=\"ocli-section\" id=\"config\"\u003e\n  \u003cdiv class=\"ocli-section__sticky\"\u003e\n    \u003cspan class=\"ocli-section__sticky-title\"\u003eConfiguration\u003c/span\u003e\n  \u003c/div\u003e\n  \u003cdiv class=\"ocli-section__body\"\u003e\n    \u003ctable class=\"ocli-table\"\u003e\n      \u003cthead\u003e\n        \u003ctr\u003e\n          \u003cth\u003eFormat\u003c/th\u003e\n          \u003cth\u003ePath\u003c/th\u003e\n        \u003c/tr\u003e\n      \u003c/thead\u003e\n      \u003ctbody\u003e\n        \u003ctr\u003e\u003ctd\u003eYAML\u003c/td\u003e\u003ctd\u003e\u003ccode\u003e~/.ocli/config.yaml\u003c/code\u003e\u003c/td\u003e\u003c/tr\u003e\n      \u003c/tbody\u003e\n    \u003c/table\u003e\n  \u003c/div\u003e\n\u003c/section\u003e\n\n      \n\u003csection class=\"ocli-section\" id=\"exit-codes\"\u003e\n  \u003cdiv class=\"ocli-section__sticky\"\u003e\n    \u003cspan class=\"ocli-section__sticky-title\"\u003eExit Codes\u003c/span\u003e\n  \u003c/div\u003e\n  \u003cdiv class=\"ocli-section__body\"\u003e\n    \u003cdiv class=\"ocli-exit-code\"\u003e\n      \u003cspan class=\"ocli-exit-code__code\"\u003e0\u003c/span\u003e\n      \u003cdiv class=\"ocli-exit-code__details\"\u003e\n        \u003cdiv class=\"ocli-exit-code__status\"\u003eOK\u003c/div\u003e\n        \u003cdiv class=\"ocli-exit-code__summary\"\u003eThe command was successful\u003c/div\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n    \u003cdiv class=\"ocli-exit-code\"\u003e\n      \u003cspan class=\"ocli-exit-code__code\"\u003e1\u003c/span\u003e\n      \u003cdiv class=\"ocli-exit-code__details\"\u003e\n        \u003cdiv class=\"ocli-exit-code__status\"\u003eINTERNAL_CLI_ERROR\u003c/div\u003e\n        \u003cdiv class=\"ocli-exit-code__summary\"\u003eThe command failed\u003c/div\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n    \u003cdiv class=\"ocli-exit-code\"\u003e\n      \u003cspan class=\"ocli-exit-code__code\"\u003e2\u003c/span\u003e\n      \u003cdiv class=\"ocli-exit-code__details\"\u003e\n        \u003cdiv class=\"ocli-exit-code__status\"\u003eBAD_USER_INPUT_ERROR\u003c/div\u003e\n        \u003cdiv class=\"ocli-exit-code__summary\"\u003eMissing or invalid input\u003c/div\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n  \u003c/div\u003e\n\u003c/section\u003e\n\n\u003csection class=\"ocli-section\" id=\"ocli-check\"\u003e\n  \u003cdiv class=\"ocli-section__sticky\"\u003e\n    \u003cdiv class=\"ocli-section__cmd-anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-check\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-section__cmd\"\u003e$ ocli check \u0026lt;path-to-spec\u0026gt; [flags]\u003c/code\u003e\u003c/div\u003e\n  \u003c/div\u003e\n  \u003cdiv class=\"ocli-section__body\"\u003e\n    \u003cp class=\"ocli-summary\"\u003eCheck an OpenCLI Spec document for errors\u003c/p\u003e\n    \u003cdiv class=\"ocli-md\"\u003e\u003cp\u003eThe \u003ccode\u003echeck\u003c/code\u003e command will validate an OpenCLI Spec document and return any\nerrors that are found.\u003c/p\u003e\n\u003c/div\u003e\n    \u003cdiv class=\"ocli-params\"\u003e\n      \u003cdiv class=\"ocli-params__heading\"\u003eArguments\u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-check-arg-path-to-spec\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-check-arg-path-to-spec\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e\u0026lt;path-to-spec\u0026gt;\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003estring\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--required\"\u003erequired\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003eThe path to the OpenCLI Spec document to check\u003c/p\u003e\u003cdiv class=\"ocli-param__desc ocli-md\"\u003e\u003cp\u003ethis command will check the validity of the given OpenCLI document\nby parsing it and checking it against the OpenCLI Specification.\u003c/p\u003e\n\u003cblockquote\u003e\n\u003cp\u003e[!TIP]\nThe format of the input file will be determined by the file extension.\u003c/p\u003e\n\u003c/blockquote\u003e\n\u003c/div\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n    \u003cdiv class=\"ocli-params\"\u003e\n      \u003cdiv class=\"ocli-params__heading\"\u003eFlags\u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-check-flag-fail-on-err\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-check-flag-fail-on-err\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e--fail-on-err\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003eboolean\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--default\"\u003edefault: true\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003eReturn a non-zero exit code if validation fails\u003c/p\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n  \u003c/div\u003e\n\u003c/section\u003e\n\n\u003csection class=\"ocli-section\" id=\"ocli-gen\"\u003e\n  \u003cdiv class=\"ocli-section__sticky\"\u003e\n    \u003cdiv class=\"ocli-section__cmd-anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-section__cmd\"\u003e$ ocli gen {command} \u0026lt;arguments\u0026gt; [flags]\u003c/code\u003e\u003c/div\u003e\n    \u003cdiv class=\"ocli-badges\"\u003e\n      \u003cspan class=\"ocli-badge ocli-badge--group\"\u003egrouping\u003c/span\u003e\n    \u003c/div\u003e\n  \u003c/div\u003e\n  \u003cdiv class=\"ocli-section__body\"\u003e\n    \u003cp class=\"ocli-summary\"\u003eCollection of commands used to generate code/docs from an OpenCLI Spec document\u003c/p\u003e\n  \u003c/div\u003e\n\u003c/section\u003e\n\u003csection class=\"ocli-section\" id=\"ocli-gen-docs\"\u003e\n  \u003cdiv class=\"ocli-section__sticky\"\u003e\n    \u003cdiv class=\"ocli-section__cmd-anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-docs\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-section__cmd\"\u003e$ ocli gen docs \u0026lt;path-to-spec\u0026gt; [flags]\u003c/code\u003e\u003c/div\u003e\n  \u003c/div\u003e\n  \u003cdiv class=\"ocli-section__body\"\u003e\n    \u003cp class=\"ocli-summary\"\u003eGenerate documentation\u003c/p\u003e\n    \u003cdiv class=\"ocli-md\"\u003e\u003cp\u003eThe \u003ccode\u003eocli gen docs\u003c/code\u003e command will generate documentation from an OpenCLI\nSpec document. You can specify the format of the documentation to be\ngenerated using the \u003ccode\u003e--format\u003c/code\u003e flag.\u003c/p\u003e\n\u003c/div\u003e\n    \u003cdiv class=\"ocli-params\"\u003e\n      \u003cdiv class=\"ocli-params__heading\"\u003eArguments\u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-gen-docs-arg-path-to-spec\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-docs-arg-path-to-spec\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e\u0026lt;path-to-spec\u0026gt;\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003estring\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--required\"\u003erequired\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003ePath to the OpenCLI Spec document\u003c/p\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n    \u003cdiv class=\"ocli-params\"\u003e\n      \u003cdiv class=\"ocli-params__heading\"\u003eFlags\u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-gen-docs-flag-format\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-docs-flag-format\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e--format\u003c/code\u003e\u003ccode class=\"ocli-param__name\"\u003e-f\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003estring\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--default\"\u003eenum\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--required\"\u003erequired\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003eThe format of the documentation to generate\u003c/p\u003e\n        \u003cdiv class=\"ocli-choices\"\u003e\n          \u003cdiv class=\"ocli-choices__label\"\u003eAccepted values\u003c/div\u003e\n          \u003cul class=\"ocli-choices__list\"\u003e\n            \u003cli class=\"ocli-choices__item\"\u003e\n              \u003ccode\u003emarkdown\u003c/code\u003e\u003cspan class=\"ocli-choices__item-desc\"\u003ePlain text using standard [Markdown](https://commonmark.org).\u003c/span\u003e\n            \u003c/li\u003e\n            \u003cli class=\"ocli-choices__item\"\u003e\n              \u003ccode\u003ehtml-page\u003c/code\u003e\u003cspan class=\"ocli-choices__item-desc\"\u003eFull standalone HTML document\u003c/span\u003e\n            \u003c/li\u003e\n            \u003cli class=\"ocli-choices__item\"\u003e\n              \u003ccode\u003ehtml-embed\u003c/code\u003e\u003cspan class=\"ocli-choices__item-desc\"\u003eAn embeddable js bundle including CSS and HTML\u003c/span\u003e\n            \u003c/li\u003e\n          \u003c/ul\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-alt-sources\"\u003e\n          \u003cdiv class=\"ocli-alt-sources__label\"\u003eAlternative sources\u003c/div\u003e\n          \u003cul class=\"ocli-alt-sources__list\"\u003e\n            \u003cli class=\"ocli-alt-sources__item\"\u003e\n              \u003ccode\u003eenv: $OCLI_DOCS_FORMAT\u003c/code\u003e\n            \u003c/li\u003e\n            \u003cli class=\"ocli-alt-sources__item\"\u003e\n              \u003ccode\u003efile: $.docs.format\u003c/code\u003e\n            \u003c/li\u003e\n          \u003c/ul\u003e\n        \u003c/div\u003e\n      \u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-gen-docs-flag-out\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-docs-flag-out\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e--out\u003c/code\u003e\u003ccode class=\"ocli-param__name\"\u003e-o\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003estring\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--required\"\u003erequired\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--default\"\u003edefault: ./docs\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003eThe path to the directory where the generated docs will be written\u003c/p\u003e\n      \u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-gen-docs-flag-no-footer\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-docs-flag-no-footer\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e--no-footer\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003eboolean\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003eDo not include the footer in the docs\u003c/p\u003e\n      \u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-gen-docs-flag-no-badge\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-docs-flag-no-badge\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e--no-badge\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003eboolean\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003eDo not include the OpenCLI badge\u003c/p\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n  \u003c/div\u003e\n\u003c/section\u003e\n\n\u003csection class=\"ocli-section\" id=\"ocli-gen-cli\"\u003e\n  \u003cdiv class=\"ocli-section__sticky\"\u003e\n    \u003cdiv class=\"ocli-section__cmd-anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-cli\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-section__cmd\"\u003e$ ocli gen cli \u0026lt;path-to-spec\u0026gt; [flags]\u003c/code\u003e\u003c/div\u003e\n  \u003c/div\u003e\n  \u003cdiv class=\"ocli-section__body\"\u003e\n    \u003cp class=\"ocli-summary\"\u003eGenerate CLI Boilerplate code\u003c/p\u003e\n    \u003cdiv class=\"ocli-md\"\u003e\u003cp\u003eThe \u003ccode\u003eocli gen cli\u003c/code\u003e command will generate the boilerplate code for\nvarious CLI frameworks from an OpenCLI Spec document.\u003c/p\u003e\n\u003c/div\u003e\n    \u003cdiv class=\"ocli-params\"\u003e\n      \u003cdiv class=\"ocli-params__heading\"\u003eArguments\u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-gen-cli-arg-path-to-spec\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-cli-arg-path-to-spec\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e\u0026lt;path-to-spec\u0026gt;\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003estring\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--required\"\u003erequired\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003eThe spec used to generate the code\u003c/p\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n    \u003cdiv class=\"ocli-params\"\u003e\n      \u003cdiv class=\"ocli-params__heading\"\u003eFlags\u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-gen-cli-flag-out\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-cli-flag-out\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e--out\u003c/code\u003e\u003ccode class=\"ocli-param__name\"\u003e-o\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003estring\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--required\"\u003erequired\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003eThe path to the directory where the generated code will be output\u003c/p\u003e\n      \u003c/div\u003e\n      \u003cdiv class=\"ocli-param\" id=\"ocli-gen-cli-flag-framework\"\u003e\n        \u003cdiv class=\"ocli-param__header\"\u003e\n          \u003cdiv class=\"ocli-param__anchor\"\u003e\u003cbutton class=\"ocli-deeplink-btn\" type=\"button\" data-href=\"#ocli-gen-cli-flag-framework\" aria-label=\"Copy link\"\u003e\u003csvg class=\"ocli-deeplink-icon\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" aria-hidden=\"true\"\u003e\u003cpath d=\"M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71\"/\u003e\u003cpath d=\"M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71\"/\u003e\u003c/svg\u003e\u003c/button\u003e\u003ccode class=\"ocli-param__name\"\u003e--framework\u003c/code\u003e\u003ccode class=\"ocli-param__name\"\u003e-f\u003c/code\u003e\u003c/div\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-badges\"\u003e\n          \u003cspan class=\"ocli-badge ocli-badge--default\"\u003estring\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--default\"\u003eenum\u003c/span\u003e\u003cspan class=\"ocli-badge ocli-badge--required\"\u003erequired\u003c/span\u003e\n        \u003c/div\u003e\u003cp class=\"ocli-param__summary\"\u003eThe framework of the CLI boilerplate to generate\u003c/p\u003e\n        \u003cdiv class=\"ocli-choices\"\u003e\n          \u003cdiv class=\"ocli-choices__label\"\u003eAccepted values\u003c/div\u003e\n          \u003cul class=\"ocli-choices__list\"\u003e\n            \u003cli class=\"ocli-choices__item\"\u003e\n              \u003ccode\u003ecobra\u003c/code\u003e\u003cspan class=\"ocli-choices__item-desc\"\u003e[Cobra](https://github.com/spf13/cobra?tab=readme-ov-file) is a popular CLI framework for Go\u003c/span\u003e\n            \u003c/li\u003e\n            \u003cli class=\"ocli-choices__item\"\u003e\n              \u003ccode\u003eyargs\u003c/code\u003e\u003cspan class=\"ocli-choices__item-desc\"\u003e[yargs](https://yargs.js.org) is a popular CLI framework for Node.js\u003c/span\u003e\n            \u003c/li\u003e\n          \u003c/ul\u003e\n        \u003c/div\u003e\n        \u003cdiv class=\"ocli-alt-sources\"\u003e\n          \u003cdiv class=\"ocli-alt-sources__label\"\u003eAlternative sources\u003c/div\u003e\n          \u003cul class=\"ocli-alt-sources__list\"\u003e\n            \u003cli class=\"ocli-alt-sources__item\"\u003e\n              \u003ccode\u003eenv: $OCLI_CLI_FRAMEWORK\u003c/code\u003e\n            \u003c/li\u003e\n            \u003cli class=\"ocli-alt-sources__item\"\u003e\n              \u003ccode\u003efile: $.cli.framework\u003c/code\u003e\n            \u003c/li\u003e\n          \u003c/ul\u003e\n        \u003c/div\u003e\n      \u003c/div\u003e\n    \u003c/div\u003e\n  \u003c/div\u003e\n\u003c/section\u003e\n\n\n    \u003cfooter class=\"ocli-footer\"\u003e\n      generated by \u003ca href=\"https://opencli.dev\"\u003eOpenCLI\u003c/a\u003e\n    \u003c/footer\u003e\n  \u003c/div\u003e\n\u003c/div\u003e\n\n";


	var DEFAULT_CONTAINER_ID = "docs";
	var THEME_STORAGE_KEY = "ocli-theme";

	function asArray(nodeList) {
		return Array.prototype.slice.call(nodeList || []);
	}

	function findById(root, id) {
		if (!id) return null;
		if (global.CSS && typeof global.CSS.escape === "function") {
			return root.querySelector("#" + global.CSS.escape(id));
		}
		var escaped = id.replace(/\\/g, "\\\\").replace(/\"/g, '\\\"');
		return root.querySelector('[id="' + escaped + '"]');
	}

	function getStoredTheme(storageKey) {
		try {
			var theme = global.localStorage.getItem(storageKey);
			return theme === "light" || theme === "dark" ? theme : null;
		} catch (_err) {
			return null;
		}
	}

	function setStoredTheme(storageKey, theme) {
		try {
			global.localStorage.setItem(storageKey, theme);
		} catch (_err) {
			// Ignore environments where localStorage is unavailable.
		}
	}

	function getSystemTheme() {
		if (global.matchMedia && global.matchMedia("(prefers-color-scheme: dark)").matches) {
			return "dark";
		}
		return "light";
	}

	function applyTheme(root, theme) {
		root.setAttribute("data-theme", theme);

		var toggle = root.querySelector("[data-theme-toggle]");
		if (!toggle) return;

		var isDark = theme === "dark";
		toggle.setAttribute("aria-pressed", String(isDark));
		toggle.setAttribute("aria-label", isDark ? "Switch to light theme" : "Switch to dark theme");
		toggle.setAttribute("title", isDark ? "Switch to light theme" : "Switch to dark theme");
	}

	function initTheme(root, opts) {
		var toggle = root.querySelector("[data-theme-toggle]");
		if (!toggle) return;

		var storageKey = opts.themeStorageKey || THEME_STORAGE_KEY;
		var mediaQuery = global.matchMedia ? global.matchMedia("(prefers-color-scheme: dark)") : null;
		var storedTheme = getStoredTheme(storageKey);
		var currentTheme = storedTheme || getSystemTheme();
		applyTheme(root, currentTheme);

		toggle.addEventListener("click", function () {
			currentTheme = currentTheme === "dark" ? "light" : "dark";
			applyTheme(root, currentTheme);
			setStoredTheme(storageKey, currentTheme);
		});

		if (!storedTheme && mediaQuery) {
			var onSystemThemeChange = function (event) {
				currentTheme = event.matches ? "dark" : "light";
				applyTheme(root, currentTheme);
			};

			if (typeof mediaQuery.addEventListener === "function") {
				mediaQuery.addEventListener("change", onSystemThemeChange);
			} else if (typeof mediaQuery.addListener === "function") {
				mediaQuery.addListener(onSystemThemeChange);
			}
		}
	}

	function initScrollSpy(root) {
		var main = root.querySelector(".ocli__main");
		if (!main) return;

		var sections = asArray(main.querySelectorAll(".ocli-section[id]"));
		var navLinks = asArray(root.querySelectorAll(".ocli__nav-link[data-section]"));
		if (!sections.length || !navLinks.length) return;

		function getActiveId() {
			var scrollTop = main.scrollTop;
			var found = null;
			for (var i = 0; i < sections.length; i++) {
				var top = sections[i].offsetTop - 80;
				if (scrollTop >= top) {
					found = sections[i].id;
				} else {
					break;
				}
			}
			return found || sections[0].id;
		}

		function updateNav() {
			var activeId = getActiveId();
			navLinks.forEach(function (link) {
				if (link.dataset.section === activeId) {
					link.classList.add("is-active");
					var parentSublist = link.closest(".ocli__nav-sublist");
					if (parentSublist) {
						var parentItem = parentSublist.closest(".ocli__nav-item");
						if (parentItem) parentItem.classList.add("is-open");
					}

					var sidebar = root.querySelector(".ocli__sidebar-body");
					if (sidebar) {
						var linkRect = link.getBoundingClientRect();
						var sidebarRect = sidebar.getBoundingClientRect();
						if (linkRect.bottom > sidebarRect.bottom || linkRect.top < sidebarRect.top) {
							link.scrollIntoView({ block: "nearest", behavior: "smooth" });
						}
					}
				} else {
					link.classList.remove("is-active");
				}
			});
		}

		main.addEventListener("scroll", updateNav, { passive: true });
		updateNav();
	}

	function initNavClick(root) {
		var navLinks = asArray(root.querySelectorAll('.ocli__nav-link[href^="#"]'));
		var main = root.querySelector(".ocli__main");
		if (!main) return;

		navLinks.forEach(function (link) {
			link.addEventListener("click", function (e) {
				var href = link.getAttribute("href");
				var id = href ? href.slice(1) : "";
				var target = findById(root, id);
				if (target) {
					e.preventDefault();
					var offset = target.offsetTop - 8;
					main.scrollTo({ top: offset, behavior: "smooth" });
				}
			});
		});
	}

	function initCollapsible(root) {
		var groupItems = asArray(root.querySelectorAll(".ocli__nav-item--group"));
		groupItems.forEach(function (item) {
			var link = item.querySelector(":scope > .ocli__nav-link");
			if (!link) return;

			link.addEventListener("click", function () {
				var href = link.getAttribute("href");
				if (href && href.charAt(0) === "#") {
					item.classList.toggle("is-open");
				}
			});
		});
	}

	function initStickyHeaders(root) {
		var main = root.querySelector(".ocli__main");
		var headers = asArray(root.querySelectorAll(".ocli-section__sticky"));
		if (!("IntersectionObserver" in global) || !main) return;

		headers.forEach(function (header) {
			if (!header.parentNode) return;

			var sentinel = global.document.createElement("div");
			sentinel.style.cssText = "height:1px;pointer-events:none;visibility:hidden;";
			header.parentNode.insertBefore(sentinel, header);

			new IntersectionObserver(
				function (entries) {
					entries.forEach(function (entry) {
						header.classList.toggle("is-stuck", !entry.isIntersecting);
					});
				},
				{ root: main, threshold: 1 }
			).observe(sentinel);
		});
	}

	function initInstallTabs(root) {
		var tabLists = asArray(root.querySelectorAll(".ocli-tabs__tablist"));
		tabLists.forEach(function (tabList) {
			tabList.addEventListener("click", function (e) {
				var btn = e.target.closest(".ocli-tabs__tab");
				if (!btn) return;

				var container = tabList.closest(".ocli-tabs");
				if (!container) return;

				asArray(container.querySelectorAll(".ocli-tabs__tab")).forEach(function (tab) {
					tab.classList.remove("is-active");
					tab.setAttribute("aria-selected", "false");
				});
				asArray(container.querySelectorAll(".ocli-tabs__panel")).forEach(function (panel) {
					panel.classList.remove("is-active");
				});

				btn.classList.add("is-active");
				btn.setAttribute("aria-selected", "true");

				var panelId = btn.dataset.tab;
				var panel = panelId ? findById(root, panelId) : null;
				if (panel) panel.classList.add("is-active");
			});
		});
	}

	function initDeepLinks(root) {
		var main = root.querySelector(".ocli__main");
		asArray(root.querySelectorAll(".ocli-deeplink-btn")).forEach(function (btn) {
			btn.addEventListener("click", function (e) {
				e.preventDefault();

				var href = btn.getAttribute("data-href") || "";
				if (!href) return;

				var url = global.location.href.split("#")[0] + href;
				if (global.navigator.clipboard && global.navigator.clipboard.writeText) {
					global.navigator.clipboard.writeText(url);
				} else {
					var ta = global.document.createElement("textarea");
					ta.value = url;
					ta.style.cssText = "position:fixed;top:-9999px;left:-9999px;";
					global.document.body.appendChild(ta);
					ta.select();
					global.document.execCommand("copy");
					global.document.body.removeChild(ta);
				}

				var id = href.slice(1);
				var target = findById(root, id);
				if (target && main) {
					var isCmd = !!btn.closest(".ocli-section__cmd-anchor");
					main.scrollTo({ top: target.offsetTop - (isCmd ? 0 : 60), behavior: "smooth" });
				}
			});
		});
	}

	function initDocs(root, opts) {
		initTheme(root, opts);
		initScrollSpy(root);
		initNavClick(root);
		initCollapsible(root);
		initStickyHeaders(root);
		initInstallTabs(root);
		initDeepLinks(root);
	}

	function resolveContainer(options) {
		if (options && options.container && options.container.nodeType === 1) {
			return options.container;
		}

		var containerId = (options && options.containerId) || DEFAULT_CONTAINER_ID;
		return global.document.getElementById(containerId);
	}

	function OcliDocs(options) {
		var opts = options || {};
		var container = resolveContainer(opts);
		if (!container) {
			throw new Error("OcliDocs: target container not found");
		}

		container.innerHTML = EMBED_MARKUP;
		var root = container.querySelector(".ocli");
		if (!root) {
			throw new Error("OcliDocs: failed to mount docs markup");
		}

		initDocs(root, opts);
		return { container: container, root: root };
	}

	global.OcliDocs = OcliDocs;

})(window);
