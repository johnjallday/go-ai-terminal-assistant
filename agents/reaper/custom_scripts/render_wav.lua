reaper.GetSetProjectInfo(0, "RENDER_SETTINGS", 0, true)
reaper.GetSetProjectInfo(0, "RENDER_SRATE", 48000, true)
reaper.GetSetProjectInfo_String(0, "RENDER_FORMAT", "wav", true)
reaper.GetSetProjectInfo_String(0, "RENDER_PATTERN", "Title", true)
reaper.Main_OnCommand(42230, 0)
