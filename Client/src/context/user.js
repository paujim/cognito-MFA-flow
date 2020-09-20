import React from 'react'
import { useAuthAPI } from "../utils/auth-api";

export const UserContext = React.createContext()

export function UserProvider(props) {

    const [user, setUser] = React.useState(useAuthAPI.getUser());

    return (
        <UserContext.Provider
            value={{
                user,
                setUser,
            }} >
            {props.children}
        </UserContext.Provider >
    );
}