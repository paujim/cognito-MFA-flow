import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';

import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';
import CircularProgress from '@material-ui/core/CircularProgress';
import TextField from '@material-ui/core/TextField';
import InputAdornment from '@material-ui/core/InputAdornment';
import IconButton from '@material-ui/core/IconButton';
import SendIcon from '@material-ui/icons/Send';

import PhonelinkLockIcon from '@material-ui/icons/PhonelinkLock';

import { UserContext } from "../context/user";
import { useAuthAPI } from "../utils/auth-api";

const useStyles = makeStyles((theme) => ({
    card: {
        minWidth: 275,
    },
    wrapper: {
        margin: theme.spacing(1),
        position: 'relative',
    },
    buttonProgress: {
        position: 'absolute',
        top: '50%',
        left: '50%',
        marginTop: -8,
        marginLeft: -4,
    },
    media: {
        height: 140,
    },
    qrCode: {
        width: "100%",
    },
    pos: {
        marginBottom: 24,
    },
    spacer: {
        marginTop: 24,
        marginBottom: 24,
    },
}));

const timeConverter = (unixTimestamp) => {
    var a = new Date(unixTimestamp * 1000);
    var months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    var year = a.getFullYear();
    var month = months[a.getMonth()];
    var date = a.getDate();
    var hour = a.getHours();
    var min = a.getMinutes();
    var sec = a.getSeconds();
    var time = date + ' ' + month + ' ' + year + ' ' + hour + ':' + min + ':' + sec;
    return time;
}

const TokenText = (props) => {
    if (props.exp) {
        return (
            <Typography variant="body2" component="div" className={props.className}>
                Token expiration date: <span color="primary">{timeConverter(props.exp)}</span>
                <br />
            Issued by: <span color="textSecondary">{props.iss}</span>
            </Typography>)
    }
    return (<Typography variant="body2" component="div" > No token found </Typography>)
}

export default function Layout(props) {
    const classes = useStyles();

    const { user } = React.useContext(UserContext)

    const [isLoading, setIsLoading] = React.useState(false);
    const [code, setCode] = React.useState();
    const [mfaData, setMfaData] = React.useState({ hasMfa: false });

    const handleCodeChange = (event) => {
        setCode(event.target.value)
    }

    const handleRegisterMfa = () => {
        setIsLoading(true)
        useAuthAPI.registerMFA(useAuthAPI.getToken())
            .then(data => {
                setIsLoading(false)
                console.log(data)
                setMfaData({ data, hasMfa: true })
                console.log("MFA Regitered Successfully")

            })
            .catch(error => {
                setIsLoading(false)
                console.log("MFA Registration Falied")
            })
    }
    const handleVerifyMfaCode = () => {
        setIsLoading(true)
        useAuthAPI.verifyMFA(useAuthAPI.getToken(), code)
            .then(data => {
                setIsLoading(false)
                console.log(data)
                setMfaData({ hasMfa: false })
                console.log("MFA Code Verified Successfully")

            })
            .catch(error => {
                setIsLoading(false)
                console.log("MFA Code Verification Falied")
            })
    }

    return (
        <div style={{ padding: 20 }}>
            <Grid container spacing={5} justify="center">
                <Grid item xs={6}  >
                    <Card className={classes.card} variant="outlined">
                        <CardContent>
                            <Typography variant="h6" component="div" >
                                {user.username}
                            </Typography>
                            <Typography variant="caption" component="div" color="textSecondary" className={classes.pos}>
                                {user.sub ? user.sub : "No ID"}
                            </Typography>
                            <TokenText iss={user.iss} exp={user.exp} />
                            {/* <Typography variant="body2" component="span">
                                <pre style={{ whiteSpace: "pre-wrap" }}>{JSON.stringify(user, null, 4)}</pre>
                            </Typography> */}
                        </CardContent>
                        <Divider />
                        {mfaData.hasMfa &&
                            <CardContent>
                                <Typography gutterBottom variant="h5" component="h2">
                                    Google Autheticator QR
                                </Typography>
                                <img className={classes.qrCode} src={`data:image/png;base64,${mfaData.data.googleAutheticator}`} alt={"QR Code"} />
                                <TextField
                                    id="input-with-icon-textfield"
                                    label="code"
                                    helperText="Send the code from Google Authenticator"
                                    fullWidth
                                    variant="outlined"
                                    value={code}
                                    onChange={handleCodeChange}
                                    InputProps={{
                                        endAdornment: (
                                            <InputAdornment position="end">
                                                <IconButton
                                                    aria-label="verify code"
                                                    onClick={handleVerifyMfaCode}
                                                    // onMouseDown={handleClickVerifyCode}
                                                    edge="end"
                                                >
                                                    <SendIcon />
                                                </IconButton>
                                            </InputAdornment>
                                        ),
                                    }}
                                />
                            </CardContent>}
                        <CardActions>
                            <div className={classes.wrapper}>
                                <Button size="small" disabled={isLoading || !useAuthAPI.isAuthenticated()} startIcon={<PhonelinkLockIcon />} onClick={handleRegisterMfa}>QR Code</Button>
                                {isLoading && <CircularProgress size={24} className={classes.buttonProgress} />}
                            </div>
                        </CardActions>

                    </Card>
                </Grid>
            </Grid>
        </div>
    );
}
