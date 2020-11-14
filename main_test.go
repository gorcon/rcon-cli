package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gorcon/rcon-cli/internal/config"
	"github.com/gorcon/rcon-cli/internal/logger"
	"github.com/gorcon/rcon-cli/internal/session"
	"github.com/stretchr/testify/assert"
)

func TestReadYamlConfig(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		expected := &config.Config{
			"default": session.Session{Address: "", Password: "", Log: "rcon-default.log"},
		}

		cfg, err := config.ReadYamlConfig("rcon.yaml")
		assert.NoError(t, err)
		assert.Equal(t, expected, &cfg)
	})

	t.Run("not exist file", func(t *testing.T) {
		var expected config.Config

		cfg, err := config.ReadYamlConfig("nonexist.yaml")
		assert.NotNil(t, err)
		assert.Equal(t, expected, cfg)
	})
}

func TestAddLog(t *testing.T) {
	logName := "tmpfile.log"

	address := "127.0.0.1:16200"
	command := "players"
	result := `Players connected (2):
-admin
-testuser`

	defer func() {
		err := os.Remove(logName)
		assert.NoError(t, err)
	}()

	// Test skip log. No logs is available.
	t.Run("skip log", func(t *testing.T) {
		err := logger.AddLog("", address, command, result)
		assert.NoError(t, err)
	})

	// Test create log file.
	t.Run("create log file", func(t *testing.T) {
		err := logger.AddLog(logName, address, command, result)
		assert.NoError(t, err)
	})

	// Test append to log file.
	t.Run("append to log file", func(t *testing.T) {
		err := logger.AddLog(logName, address, command, result)
		assert.NoError(t, err)
	})
}

func TestGetLogFile(t *testing.T) {
	logDir := "temp"
	logName := "tmpfile.log"
	logPath := logDir + "/" + logName

	// Test empty log file name.
	t.Run("empty file name", func(t *testing.T) {
		file, err := logger.GetLogFile("")
		assert.Nil(t, file)
		assert.EqualError(t, err, "empty file name")
	})

	// Test stat permission denied.
	t.Run("stat permission denied", func(t *testing.T) {
		if err := os.Mkdir(logDir, 0400); err != nil {
			assert.NoError(t, err)
			return
		}
		defer func() {
			err := os.RemoveAll(logDir)
			assert.NoError(t, err)
		}()

		file, err := logger.GetLogFile(logPath)
		assert.Nil(t, file)
		assert.EqualError(t, err, fmt.Sprintf("stat %s: permission denied", logPath))
	})

	// Test create permission denied.
	t.Run("open permission denied", func(t *testing.T) {
		if err := os.Mkdir(logDir, 0500); err != nil {
			assert.NoError(t, err)
			return
		}
		defer func() {
			err := os.RemoveAll(logDir)
			assert.NoError(t, err)
		}()

		file, err := logger.GetLogFile(logPath)
		assert.Nil(t, file)
		assert.EqualError(t, err, fmt.Sprintf("open %s: permission denied", logPath))
	})

	// Positive test create new log file + test open permission denied.
	t.Run("create new log file", func(t *testing.T) {
		if err := os.Mkdir(logDir, 0700); err != nil {
			assert.NoError(t, err)
			return
		}
		defer func() {
			err := os.RemoveAll(logDir)
			assert.NoError(t, err)
		}()

		// Positive test create new log file.
		file, err := logger.GetLogFile(logPath)
		assert.NotNil(t, file)
		assert.NoError(t, err)

		if err := os.Chmod(logPath, 0000); err != nil {
			assert.NoError(t, err)
			return
		}

		// Test open permission denied.
		file, err = logger.GetLogFile(logPath)
		assert.Nil(t, file)
		assert.EqualError(t, err, fmt.Sprintf("open %s: permission denied", logPath))
	})
}

func TestExecute(t *testing.T) {
	serverRCON, err := NewMockServerRCON()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, serverRCON.Close())
		close(serverRCON.errors)
		for err := range serverRCON.errors {
			assert.NoError(t, err)
		}
	}()

	serverTELNET, err := NewMockServerTELNET()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, serverTELNET.Close())
		close(serverTELNET.errors)
		for err := range serverTELNET.errors {
			assert.NoError(t, err)
		}
	}()

	w := &bytes.Buffer{}

	// Test empty address.
	t.Run("empty address", func(t *testing.T) {
		err := Execute(w, session.Session{Address: "", Password: MockPasswordRCON}, MockCommandHelpRCON)
		assert.Error(t, err)
	})

	// Test empty password.
	t.Run("empty password", func(t *testing.T) {
		err := Execute(w, session.Session{Address: serverRCON.Addr(), Password: ""}, MockCommandHelpRCON)
		assert.Error(t, err)
	})

	// Test wrong password.
	t.Run("wrong password", func(t *testing.T) {
		err := Execute(w, session.Session{Address: serverRCON.Addr(), Password: "wrong"}, MockCommandHelpRCON)
		assert.Error(t, err)
	})

	// Test empty command.
	t.Run("empty command", func(t *testing.T) {
		err := Execute(w, session.Session{Address: serverRCON.Addr(), Password: MockPasswordRCON}, "")
		assert.Error(t, err)
	})

	// Test long command.
	t.Run("long command", func(t *testing.T) {
		bigCommand := make([]byte, 1001)
		err := Execute(w, session.Session{Address: serverRCON.Addr(), Password: MockPasswordRCON}, string(bigCommand))
		assert.Error(t, err)
	})

	// Positive RCON test Execute func.
	t.Run("no error rcon", func(t *testing.T) {
		err := Execute(w, session.Session{Address: serverRCON.Addr(), Password: MockPasswordRCON}, MockCommandHelpRCON)
		assert.NoError(t, err)
	})

	// Positive TELNET test Execute func.
	t.Run("no error telnet", func(t *testing.T) {
		err := Execute(w, session.Session{Address: serverTELNET.Addr(), Password: MockPasswordTELNET, Type: session.ProtocolTELNET}, MockCommandHelpTELNET)
		assert.NoError(t, err)
	})

	// Positive test Execute func with log.
	t.Run("no error with log", func(t *testing.T) {
		logFileName := "tmpfile.log"
		defer func() {
			err := os.Remove(logFileName)
			assert.NoError(t, err)
		}()

		err := Execute(w, session.Session{Address: serverRCON.Addr(), Password: MockPasswordRCON, Log: logFileName}, MockCommandHelpRCON)
		assert.NoError(t, err)
	})

	if run := getVar("TEST_PZ_SERVER", "false"); run == "true" {
		addr := getVar("TEST_PZ_SERVER_ADDR", "127.0.0.1:16260")
		password := getVar("TEST_PZ_SERVER_PASSWORD", "docker")

		t.Run("pz server", func(t *testing.T) {
			needle := `List of server commands :
* addalltowhitelist : Add all the current users connected with a password in the whitelist, so their account is protected.
* additem : Add an item to a player, if no username is given the item will be added to you, count is optional, use /additem \"username\" \"module.item\" count, ex : /additem \"rj\" \"Base.Axe\" count
* adduser : Use this command to add a new user in a whitelisted server, use : /adduser \"username\" \"pwd\"
* addusertowhitelist : Add the user connected with a password in the whitelist, so his account is protected, use : /addusertowhitelist \"username\"
* addvehicle : Spawn a new vehicle, use: /addvehicle \"script\" \"user or x,y,z\", ex /addvehicle \"Base.VanAmbulance\" \"rj\"
* addxp : Add experience points to a player, use : /addxp \"playername\" perkname=xp, ex /addxp \"rj\" Woodwork=2
* alarm : Sound a building alarm at the admin's position.  Must be in a room.
* banid : Ban a SteamID, use : /banid SteamID
* banuser : Ban a user, add a -ip to also ban his ip, add a -r \"reason\" to specify a reason for the ban, use : /banuser \"username\" -ip -r \"reason\", ex /banuser \"rj\" -ip -r \"spawn kill\"
* changeoption : Use this to change a server option, use : /changeoption optionName \"newValue\"
* chopper : Start the choppers (do noise on a random player)
* createhorde : Use this to spawn a horde near a player, use : /createhorde count \"username\", ex /createhorde 150 \"rj\", username is optional except from the server console.
* godmod : Set a player invincible, if no username set it make you invincible, if no value it toggle it, use : /godmode \"username\" -value, ex /godmode \"rj\" -true (could be -false)
* gunshot : Start a gunshot (do noise on a random player)
* help : Help
* invisible : Set a player invisible zombie will ignore him, if no username set it make you invisible, if no value it toggle it, use : /invisible \"username\" -value, ex /invisible \"rj\" -true (could be -false)
* kickuser : Kick a user, add a -r \"reason\" to specify a reason for the kick, use : /kickuser \"username\" -r \"reason\"
* noclip : A player with noclip won't collide on anything, if no value it toggle it, use : /noclip \"username\" -value, ex /noclip \"rj\" -true (could be -false)
* players : List the players connected
* quit : Quit the server (but save it before)
* releasesafehouse : Release a safehouse you are the owner of, use : /releasesafehouse
* reloadlua : Reload a Lua script, use : /reloadlua \"filename\"
* reloadoptions : Reload the options on the server (ServerOptions.ini) and send them to the clients
* removeuserfromwhitelist : Remove the user from the whitelist, use: /removeuserfromwhitelist \"username\"
* save : Save the current world
* sendpulse : Toggle sending server performance info to this client, use : /sendpulse
* servermsg : Use this to broadcast a message to all connected players, use : /servermsg my message !
* setaccesslevel : Use it to set new access level to a player, acces level: admin, moderator, overseer, gm, observer. use : /setaccesslevel \"username\" \"accesslevel\", ex: /setaccesslevel \"rj\" \"moderator\"
* showoptions : Show the list of current Server options with their values.
* startrain : Start rain on the server
* stoprain : Stop rain on the server
* teleport : Teleport to a player, once teleported, wait 2 seconds to show map, use : /teleport \"playername\" or /teleport \"player1\" \"player2\", ex /teleport \"rj\" or /teleport \"rj\" \"toUser\"
* teleportto : Teleport to coordinates, use: /teleportto x,y,z, ex /teleportto 100098,189980,0
* unbanid : Unban a SteamID, use : /unbanid SteamID
* unbanuser : Unban a player, use : /unbanuser \"username\"
* voiceban : Block voice from user \"username\", use : /voiceban \"username\" -value, ex /voiceban \"rj\" -true (could be -false)`

			needle = strings.Replace(needle, "List of server commands :", "List of server commands : ", -1)

			err := Execute(w, session.Session{Address: addr, Password: password}, "help")
			assert.NoError(t, err)
			assert.NotEmpty(t, w.String())

			if !strings.Contains(w.String(), needle) {
				diff := struct {
					R string
					N string
				}{R: w.String(), N: needle}

				js, _ := json.Marshal(diff)
				fmt.Println(string(js))

				t.Error("response is not contain needle string")
			}
		})
	}

	if run := getVar("TEST_7DTD_SERVER", "false"); run == "true" {
		addr := getVar("TEST_7DTD_SERVER_ADDR", "172.22.0.2:8081")
		password := getVar("TEST_7DTD_SERVER_PASSWORD", "banana")

		t.Run("7dtd server", func(t *testing.T) {
			needle := `*** List of Commands ***
 admin => Manage user permission levels
 aiddebug => Toggles AIDirector debug output.
 audio => Watch audio stats
 automove => Player auto movement
 ban => Manage ban entries
 bents => Switches block entities on/off
 BiomeParticles => Debug
 buff => Applies a buff to the local player
 buffplayer => Apply a buff to a player
 chunkcache cc => shows all loaded chunks in cache
 chunkobserver co => Place a chunk observer on a given position.
 chunkreset cr => resets the specified chunks
 commandpermission cp => Manage command permission levels
 creativemenu cm => enables/disables the creativemenu
 debuff => Removes a buff from the local player
 debuffplayer => Remove a buff from a player
 debugmenu dm => enables/disables the debugmenu ` + `
 debugshot dbs => Lets you make a screenshot that will have some generic info
on it and a custom text you can enter. Also stores a list
of your current perk levels in a CSV file next to it.
 debugweather => Dumps internal weather state to the console.
 decomgr => ` + `
 dms => Gives control over Dynamic Music functionality.
 dof => Control DOF
 enablescope es => toggle debug scope
 exhausted => Makes the player exhausted.
 exportcurrentconfigs => Exports the current game config XMLs
 exportprefab => Exports a prefab from a world area
 floatingorigin fo => ` + `
 fov => Camera field of view
 gamestage => usage: gamestage - displays the gamestage of the local player.
 getgamepref gg => Gets game preferences
 getgamestat ggs => Gets game stats
 getoptions => Gets game options
 gettime gt => Get the current game time
 gfx => Graphics commands
 givequest => usage: givequest questname
 giveself => usage: giveself itemName [qualityLevel=6] [count=1] [putInInventory=false] [spawnWithMods=true]
 giveselfxp => usage: giveselfxp 10000
 help => Help on console and specific commands
 kick => Kicks user with optional reason. "kick playername reason"
 kickall => Kicks all users with optional reason. "kickall reason"
 kill => Kill a given entity
 killall => Kill all entities
 lgo listgameobjects => List all active game objects
 lights => Debug views to optimize lights
 listents le => lists all entities
 listplayerids lpi => Lists all players with their IDs for ingame commands
 listplayers lp => lists all players
 listthreads lt => lists all threads
 loggamestate lgs => Log the current state of the game
 loglevel => Telnet/Web only: Select which types of log messages are shown
 mem => Prints memory information and unloads resources or changes garbage collector
 memcl => Prints memory information on client and calls garbage collector
 occlusion => Control OcclusionManager
 pirs => tbd
 pois => Switches distant POIs on/off
 pplist => Lists all PersistentPlayer data
 prefab => ` + `
 prefabupdater => ` + `
 profilenetwork => Writes network profiling information
 profiling => Enable Unity profiling for 300 frames
 removequest => usage: removequest questname
 repairchunkdensity rcd => check and optionally fix densities of a chunk
 saveworld sa => Saves the world manually.
 say => Sends a message to all connected clients
 ScreenEffect => Sets a screen effect
 setgamepref sg => sets a game pref
 setgamestat sgs => sets a game stat
 settargetfps => Set the target FPS the game should run at (upper limit)
 settempunit stu => Set the current temperature units.
 settime st => Set the current game time
 show => Shows custom layers of rendering.
 showalbedo albedo => enables/disables display of albedo in gBuffer
 showchunkdata sc => shows some date of the current chunk
 showClouds => Artist command to show one layer of clouds.
 showhits => Show hit entity locations
 shownexthordetime => Displays the wandering horde time
 shownormals norms => enables/disables display of normal maps in gBuffer
 showspecular spec => enables/disables display of specular values in gBuffer
 showswings => Show melee swing arc rays
 shutdown => shuts down the game
 sleeper => Show sleeper info
 smoothworldall swa => Applies some batched smoothing commands.
 sounddebug => Toggles SoundManager debug output.
 spawnairdrop => Spawns an air drop
 spawnentity se => spawns an entity
 spawnentityat sea => Spawns an entity at a give position
 spawnscouts => Spawns zombie scouts
 SpawnScreen => Display SpawnScreen
 spawnsupplycrate => Spawns a supply crate where the player is
 spawnwanderinghorde spawnwh => Spawns a wandering horde of zombies
 spectator spectatormode sm => enables/disables spectator mode
 spectrum => Force a particular lighting spectrum.
 stab => stability
 starve hungry food => Makes the player starve (optionally specify the amount of food you want to have in percent).
 switchview sv => Switch between fpv and tpv
 SystemInfo => List SystemInfo
 teleport tp => Teleport the local player
 teleportplayer tele => Teleport a given player
 thirsty water => Makes the player thirsty (optionally specify the amount of water you want to have in percent).
 traderarea => ...
 trees => Switches trees on/off
 updatelighton => Commands for UpdateLightOnAllMaterials and UpdateLightOnPlayers
 version => Get the currently running version of the game and loaded mods
 visitmap => Visit an given area of the map. Optionally run the density check on each visited chunk.
 water => Control water settings
 weather => Control weather settings
 weathersurvival => Enables/disables weather survival
 whitelist => Manage whitelist entries
 wsmats workstationmaterials => Set material counts on workstations.
 xui => Execute XUi operations
 xuireload => Access xui related functions such as reinitializing a window group, opening a window group
 zip => Control zipline settings`

			needle = strings.Replace(needle, "\n", "\r\n", -1)
			needle = strings.Replace(needle, "some generic info\r\n", "some generic info\n", -1)
			needle = strings.Replace(needle, "Also stores a list\r\n", "Also stores a list\n", -1)

			err := Execute(w, session.Session{Address: addr, Password: password, Type: session.ProtocolTELNET}, "help")
			assert.NoError(t, err)
			assert.NotEmpty(t, w.String())

			if !strings.Contains(w.String(), needle) {
				diff := struct {
					R string
					N string
				}{R: w.String(), N: needle}

				js, _ := json.Marshal(diff)
				fmt.Println(string(js))

				t.Error("response is not contain needle string")
			}
		})
	}

	if run := getVar("TEST_RUST_SERVER_RCON", "false"); run == "true" {
		addr := getVar("TEST_RUST_SERVER_RCON_ADDR", "127.0.0.1:28016")
		password := getVar("TEST_RUST_SERVER_RCON_PASSWORD", "docker")

		t.Run("rust server rcon", func(t *testing.T) {
			err := Execute(w, session.Session{Address: addr, Password: password}, "status")
			assert.NoError(t, err)
			assert.NotEmpty(t, w.String())

			fmt.Println(w.String())
		})
	}

	if run := getVar("TEST_RUST_SERVER_WEB", "false"); run == "true" {
		addr := getVar("TEST_RUST_SERVER_WEB_ADDR", "127.0.0.1:28016")
		password := getVar("TEST_RUST_SERVER_WEB_PASSWORD", "docker")

		t.Run("rust server web", func(t *testing.T) {
			err := Execute(w, session.Session{Address: addr, Password: password, Type: session.ProtocolWebRCON}, "status")
			assert.NoError(t, err)
			assert.NotEmpty(t, w.String())

			fmt.Println(w.String())
		})
	}
}

func TestInteractive(t *testing.T) {
	server, err := NewMockServerRCON()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, server.Close())
		close(server.errors)
		for err := range server.errors {
			assert.NoError(t, err)
		}
	}()

	w := &bytes.Buffer{}

	// Test wrong password.
	t.Run("wrong password", func(t *testing.T) {
		var r bytes.Buffer
		r.WriteString(CommandQuit + "\n")

		err = Interactive(&r, w, session.Session{Address: server.Addr(), Password: "fake"})
		assert.Error(t, err)
	})

	// Test get Interactive address.
	t.Run("interactive get address", func(t *testing.T) {
		var r bytes.Buffer
		r.WriteString(server.Addr() + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(&r, w, session.Session{Address: "", Password: MockPasswordRCON})
		assert.NoError(t, err)
	})

	// Test get Interactive password.
	t.Run("interactive get password", func(t *testing.T) {
		var r bytes.Buffer
		r.WriteString(MockPasswordRCON + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(&r, w, session.Session{Address: server.Addr(), Password: ""})
		assert.NoError(t, err)
	})

	// Test get Interactive commands RCON.
	t.Run("interactive get commands rcon", func(t *testing.T) {
		r := &bytes.Buffer{}
		r.WriteString(MockCommandHelpRCON + "\n")
		r.WriteString("unknown command" + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(r, w, session.Session{Address: server.Addr(), Password: MockPasswordRCON})
		assert.NoError(t, err)
	})

	// Test get Interactive commands TELNET.
	t.Run("interactive get commands telnet", func(t *testing.T) {
		r := &bytes.Buffer{}
		r.WriteString(MockCommandHelpTELNET + "\n")
		r.WriteString("unknown command" + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(r, w, session.Session{Address: server.Addr(), Password: MockPasswordTELNET, Type: session.ProtocolTELNET})
		assert.NoError(t, err)
	})
}

func TestNewApp(t *testing.T) {
	server, err := NewMockServerRCON()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, server.Close())
		close(server.errors)
		for err := range server.errors {
			assert.NoError(t, err)
		}
	}()

	// Test getting address and password from args. Config ang log are not used.
	t.Run("getting address and password from args", func(t *testing.T) {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-a="+server.Addr())
		args = append(args, "-p="+MockPasswordRCON)
		args = append(args, "-c="+MockCommandHelpRCON)

		err = app.Run(args)
		assert.NoError(t, err)
	})

	// Test getting address and password from config. Log is not used.
	t.Run("getting address and password from args with log", func(t *testing.T) {
		var configFileName = "rcon-temp.yaml"
		err := CreateConfigFile(configFileName, server.Addr(), MockPasswordRCON)
		assert.NoError(t, err)
		defer func() {
			err := os.Remove(configFileName)
			assert.NoError(t, err)
		}()

		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-cfg="+configFileName)
		args = append(args, "-c="+MockCommandHelpRCON)

		err = app.Run(args)
		assert.NoError(t, err)
	})

	// Test default config file not exist. Log is not used.
	t.Run("default config file not exist", func(t *testing.T) {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-c="+MockCommandHelpRCON)

		err = app.Run(args)
		assert.Error(t, err)
		if !os.IsNotExist(err) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// Test default config file is incorrect. Log is not used.
	t.Run("default config file is incorrect", func(t *testing.T) {
		var configFileName = "rcon-temp.yaml"
		err := CreateInvalidConfigFile(configFileName, server.Addr(), MockPasswordRCON)
		assert.NoError(t, err)
		defer func() {
			err := os.Remove(configFileName)
			assert.NoError(t, err)
		}()

		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-cfg="+configFileName)
		args = append(args, "-c="+MockCommandHelpRCON)

		err = app.Run(args)
		assert.EqualError(t, err, "yaml: line 1: did not find expected key")
	})

	// Test empty address and password. Log is not used.
	t.Run("empty address and password", func(t *testing.T) {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		// Hack to use os.Args[0] in go run
		args[0] = ""
		args = append(args, "-c="+MockCommandHelpRCON)

		err = app.Run(args)
		assert.EqualError(t, err, "address is not set: to set address add -a host:port")
	})

	// Test empty password. Log is not used.
	t.Run("empty password", func(t *testing.T) {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		// Hack to use os.Args[0] in go run
		args[0] = ""
		args = append(args, "-a="+server.Addr())
		args = append(args, "-c="+MockCommandHelpRCON)

		err = app.Run(args)
		assert.EqualError(t, err, "password is not set: to set password add -p password")
	})

	// Positive test Interactive. Log is not used.
	t.Run("no error", func(t *testing.T) {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-a="+server.Addr())
		args = append(args, "-p="+MockPasswordRCON)

		r.WriteString(MockCommandHelpRCON + "\n")
		r.WriteString(CommandQuit + "\n")

		err = app.Run(args)
		assert.NoError(t, err)
	})
}

// DefaultTestLogName sets the default log file name.
const DefaultTestLogName = "rcon-default.log"

// CreateConfigFile creates config file with default section.
func CreateConfigFile(name string, address string, password string) error {
	var stringBody = fmt.Sprintf(
		"%s:\n  address: \"%s\"\n  password: \"%s\"\n  log: \"%s\"",
		config.DefaultConfigEnv, address, password, DefaultTestLogName,
	)
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = file.WriteString(stringBody)

	return err
}

// CreateIncorrectConfigFile creates incorrect yaml config file.
func CreateInvalidConfigFile(name string, address string, password string) error {
	var stringBody = fmt.Sprintf(
		"address: \"%s\"\n  password: \"%s\"\n  log: \"%s\"",
		address, password, DefaultTestLogName,
	)
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = file.WriteString(stringBody)

	return err
}

// getVar returns environment variable or default value.
func getVar(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
