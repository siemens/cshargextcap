use framework "Foundation"
use scripting additions

on open location thisURL
	set theURL to current application's NSURL's URLWithString:thisURL
	
	set theComponents to current application's NSURLComponents's componentsWithURL:theURL resolvingAgainstBaseURL:false
	set queryItems to theComponents's queryItems
	set containerValue to ""
	
	set theEnumerator to queryItems's objectEnumerator()
	set aQueryItem to theEnumerator's nextObject()
	
	repeat while aQueryItem is not missing value
		if ((aQueryItem's valueForKey:"name") as string is "container") then
			set containerValue to ((aQueryItem's valueForKey:"value") as string)
			exit repeat
		end if
		set aQueryItem to theEnumerator's nextObject()
	end repeat
	-- uncomment to see the parsed json from the URL
	--	display dialog (containerValue as string)
	
	
	try
		-- Convert the containerValue to NSData
		set nsContainerValue to current application's NSString's stringWithString:containerValue
		set theData to nsContainerValue's dataUsingEncoding:(current application's NSUTF8StringEncoding)
		
		-- Parse the containerValue as JSON
		set theDict to current application's NSJSONSerialization's JSONObjectWithData:theData options:0 |error|:(missing value)
		
		-- Get the value of the namespace/container "name" key
		set nodeName to (theDict's valueForKey:"name") as string
		
		-- Get the value of the "network-interfaces" key
		set networkInterfacesList to (theDict's valueForKey:"network-interfaces") as list
		
		-- Concatenate the elements into a comma-separated string
		set netIFs to ""
		repeat with i from 1 to count networkInterfacesList
			set netIFs to netIFs & item i of networkInterfacesList
			if i is not (count networkInterfacesList) then
				set netIFs to netIFs & ", "
			end if
		end repeat
		
		-- Display the value in a dialog box
		--		display dialog (netIFs as string)
		--		display dialog (nodeName as string)
		-- open wireshark
		-- create a title for wireshark window
		set title to nodeName & ": " & netIFs
	on error errMsg number errNum
		display dialog ("Error " & errNum & ": " & errMsg)
	end try
	do shell script "open -n -a /Applications/Wireshark.app/Contents/MacOS/Wireshark --args -k -i packetflix -o gui.window_title:\"" & title & "\" -o extcap.packetflix.url:\"" & thisURL & "\""
end open location
