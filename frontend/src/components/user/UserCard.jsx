import { useNavigate } from "react-router"
import avatar from "../../assets/user.svg"
import "../../styles/user-card.css"

const toDisplayName = (user) => `${user.first_name || ""} ${user.last_name || ""}`.trim() || user.nickname || "user"

const UserCard = ({user, actions, onClick}) => {
    const navigate = useNavigate()

    const handleClick = () => {
        if (onClick) {
            onClick(user)
        } else {
            navigate(`/user/${user.id}`)
        }
    }

    return (
        <div className="user-card" onClick={handleClick}>
            <img src={user.avatar|| avatar} alt={`${toDisplayName(user)}'s avatar`} className="user-card__avatar" />
            <div className="user-card__info">
                <div className="user-card__name">{toDisplayName(user)}</div>
                {user.nickname && <div className="user-card__nickname">@{user.nickname}</div>}
            </div>
        </div>
    )
}