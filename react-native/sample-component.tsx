import * as React from "react"
import { View, ViewStyle, TextStyle, SafeAreaView } from "react-native"
import { NavigationInjectedProps } from "react-navigation"
import VideoRecorder from "react-native-beautiful-video-recorder"
import * as ARKit from "@dashtracker/react-native-arkit"

import { Button, Header, Screen, Text, Wallpaper } from "../../components"
import { useStores } from "../../models/root-store"
import { color, spacing } from "../../theme"
import { useNavigation } from "@react-navigation/native"

const FULL: ViewStyle = { flex: 1 }
const CONTAINER: ViewStyle = {
  backgroundColor: color.transparent,
  paddingHorizontal: spacing[4],
}
const TEXT: TextStyle = {
  color: color.palette.white,
  fontFamily: "Montserrat",
}
const BOLD: TextStyle = { fontWeight: "bold" }
const HEADER: TextStyle = {
  paddingTop: spacing[3],
  paddingBottom: spacing[4] + spacing[1],
  paddingHorizontal: 0,
}
const HEADER_TITLE: TextStyle = {
  ...TEXT,
  ...BOLD,
  fontSize: 12,
  lineHeight: 15,
  textAlign: "center",
  letterSpacing: 1.5,
}
const TITLE_WRAPPER: TextStyle = {
  ...TEXT,
  textAlign: "center",
}
const TITLE: TextStyle = {
  ...TEXT,
  ...BOLD,
  fontSize: 20,
  lineHeight: 38,
  textAlign: "center",
}
const STEPS: TextStyle = {
  ...TITLE,
  fontSize: 16,
  textAlign: "left",
}
const CONTINUE_TEXT: TextStyle = {
  ...TEXT,
  ...BOLD,
  fontSize: 13,
  letterSpacing: 2,
}
const FOOTER: ViewStyle = { backgroundColor: "#20162D" }
const FOOTER_CONTENT: ViewStyle = {
  alignItems: "center",
  paddingVertical: spacing[4],
  paddingHorizontal: spacing[4],
  width: "100%",
}

const btnCaptureStyle: ViewStyle = {
  width: 200,
  height: 50,
  alignItems: "center",
  justifyContent: "center",
  backgroundColor: "blue",
  borderRadius: 25,
}

export interface MapThePlaneScreenProps extends NavigationInjectedProps<{}> {}

export const MapThePlaneScreen: React.FunctionComponent<MapThePlaneScreenProps> = props => {
  const navigation = useNavigation()

  const {
    permissions: { requestCameraPermission, requestMicrophonePermission },
    session: { videos },
  } = useStores()

  const videoRecorder: VideoRecorder = React.useRef(null)
  const startRecorder = () => {
    if (videoRecorder?.current) {
      videoRecorder.current.open({ maxLength: 30 }, data => {
        /*
          The calback will be fulfilled with an object with some of the following properties:

          uri: (string) the path to the video saved on your app's cache directory.
          videoOrientation: (number) orientation of the video
          deviceOrientation: (number) orientation of the device
          iOS codec: the codec of the recorded video. One of RNCamera.Constants.VideoCodec
          isRecordingInterrupted: (boolean) whether the app has been minimized while recording
        */
        console.log("captured data", data)
        videos.pushUnique(data)
        navigation.navigate('SessionVideoList')
      })
    }
  }

  let isInstalled = false
  let isSupported = false

  React.useEffect(() => {
    const requestPermissions = async () => {
      await requestCameraPermission()
      await requestMicrophonePermission()
      isSupported = await ARKit.isSupported()
      console.log("is ARKit supported?", isSupported)
      if (isSupported) {
        isInstalled = await ARKit.requestInstall()
        console.log("is ARKit installed?", isInstalled)
      }
    }

    requestPermissions()

    return () => {
      return null
    }
  }, [])

  const wallpaper = require("./wallpaper-40-yard.png")

  return (
    <View testID="MapThePlaneScreen" style={FULL}>
      <Wallpaper backgroundImage={wallpaper} />
      <Screen style={CONTAINER} preset="scroll" backgroundColor="#00000050">
        <Header
          headerTx="mapThePlaneScreen.instructions"
          style={HEADER}
          titleStyle={HEADER_TITLE}
        />
        <Text style={TITLE_WRAPPER}>
          <Text style={TITLE} text="Before you record your dashes, " />
          <Text style={TITLE} text="you will need to setup your running plane." />
          <Text style={TITLE} text={"\n\n"} />
          <Text style={STEPS} text="1. Face the camera towards the ground." />
          <Text style={TITLE} text={"\n"} />
          <Text style={STEPS} text="2. Start the mapping recorder." />
          <Text style={TITLE} text={"\n"} />
          <Text style={STEPS} text="3. Walk from your start to the end." />
          <Text style={TITLE} text={"\n\n"} />
          <Text style={TITLE} text="When you are ready, tap the button below." />
        </Text>
      </Screen>
      <SafeAreaView style={FOOTER}>
        <View style={FOOTER_CONTENT}>
          <Button
            testID="next-screen-button"
            style={btnCaptureStyle}
            textStyle={CONTINUE_TEXT}
            tx="mapThePlaneScreen.getStarted"
            onPress={startRecorder}
          />
        </View>
        <VideoRecorder ref={videoRecorder} compressQuality={"medium"} />
      </SafeAreaView>
    </View>
  )
}
