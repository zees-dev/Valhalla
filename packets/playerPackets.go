package packets

import (
	"crypto/rand"
	"fmt"
	"math"

	"github.com/Hucaru/Valhalla/consts/opcodes"
	"github.com/Hucaru/Valhalla/maplepacket"
	"github.com/Hucaru/Valhalla/nx"
	"github.com/Hucaru/Valhalla/types"
)

func PlayerReceivedDmg(charID int32, ammount int32, dmgType byte, mobID int32, hit byte, reduction byte, stance byte) maplepacket.Packet {
	p := maplepacket.CreateWithOpcode(opcodes.Send.ChannelPlayerTakeDmg)
	p.WriteInt32(charID)
	p.WriteByte(dmgType)

	if dmgType == 0xFE {
		p.WriteInt32(ammount)
		p.WriteInt32(ammount)
	} else {
		p.WriteInt32(0) // ?
		p.WriteInt32(mobID)
		p.WriteByte(hit)
		p.WriteByte(stance)
		p.WriteInt32(0)       // ?
		p.WriteInt32(ammount) // skill id of attack?
	}

	return p
}

func PlayerLevelUpAnimation(charID int32) maplepacket.Packet {
	p := maplepacket.CreateWithOpcode(opcodes.Send.ChannelPlayerAnimation)
	p.WriteInt32(charID)
	p.WriteByte(0x00)

	return p
}

func PlayerMove(charID int32, bytes []byte) maplepacket.Packet {
	p := maplepacket.CreateWithOpcode(opcodes.Send.ChannelPlayerMovement)
	p.WriteInt32(charID)
	p.WriteBytes(bytes)

	return p
}

func PlayerEmoticon(playerID int32, emotion int32) maplepacket.Packet {
	p := maplepacket.CreateWithOpcode(opcodes.Send.ChannelPlayerEmoticon)
	p.WriteInt32(playerID)
	p.WriteInt32(emotion)

	return p
}

func PlayerSkillBookUpdate(skillID int32, level int32) maplepacket.Packet {
	p := maplepacket.CreateWithOpcode(opcodes.Send.ChannelSkillRecordUpdate)
	p.WriteByte(0x01)  // time check?
	p.WriteInt16(0x01) // number of skills to update
	p.WriteInt32(skillID)
	p.WriteInt32(level)
	p.WriteByte(0x01)

	return p
}

func PlayerStatChange(byPlayer bool, stat int32, value int32) maplepacket.Packet {
	p := maplepacket.CreateWithOpcode(opcodes.Send.ChannelStatChange)
	p.WriteBool(byPlayer)
	p.WriteInt32(stat)
	p.WriteInt32(value)

	return p
}

func PlayerStatNoChange() maplepacket.Packet {
	p := maplepacket.CreateWithOpcode(opcodes.Send.ChannelInventoryOperation)
	p.WriteByte(0x01)
	p.WriteByte(0x00)
	p.WriteByte(0x00)

	return p
}

func PlayerAvatarSummaryWindow(charID int32, char types.Character, guildName string) maplepacket.Packet {
	p := maplepacket.CreateWithOpcode(opcodes.Send.ChannelAvatarInfoWindow)
	p.WriteInt32(charID)
	p.WriteByte(char.Level)
	p.WriteInt16(char.Job)
	p.WriteInt16(char.Fame)

	p.WriteString(guildName)

	p.WriteBool(false) // if has pet
	p.WriteByte(0)     // wishlist count

	return p
}

func PlayerEnterGame(char types.Character, channelID int32) maplepacket.Packet {
	p := maplepacket.CreateWithOpcode(opcodes.Send.ChannelWarpToMap)
	p.WriteInt32(channelID)
	p.WriteByte(0) // character portal counter
	p.WriteByte(1) // Is connecting

	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err.Error())
	}
	p.WriteBytes(randomBytes)
	p.WriteBytes(randomBytes)
	p.WriteBytes(randomBytes)
	p.WriteBytes(randomBytes)
	p.WriteBytes([]byte{0xFF, 0xFF}) // seperators? For what?
	p.WriteInt32(char.ID)
	p.WritePaddedString(char.Name, 13)
	p.WriteByte(char.Gender)
	p.WriteByte(char.Skin)
	p.WriteInt32(char.Face)
	p.WriteInt32(char.Hair)

	p.WriteInt64(0) // Pet Cash ID

	p.WriteByte(char.Level)
	p.WriteInt16(char.Job)
	p.WriteInt16(char.Str)
	p.WriteInt16(char.Dex)
	p.WriteInt16(char.Int)
	p.WriteInt16(char.Luk)
	p.WriteInt16(char.HP)
	p.WriteInt16(char.MaxHP)
	p.WriteInt16(char.MP)
	p.WriteInt16(char.MaxMP)
	p.WriteInt16(char.AP)
	p.WriteInt16(char.SP)
	p.WriteInt32(char.EXP)
	p.WriteInt16(char.Fame)

	p.WriteInt32(char.CurrentMap)
	p.WriteByte(char.CurrentMapPos)

	p.WriteByte(20) // budy list size
	p.WriteInt32(char.Mesos)

	p.WriteByte(char.EquipSlotSize)
	p.WriteByte(char.UseSlotSize)
	p.WriteByte(char.SetupSlotSize)
	p.WriteByte(char.EtcSlotSize)
	p.WriteByte(char.CashSlotSize)

	for _, v := range char.Equip {
		if v.SlotID < 0 && v.InvID == 1 && !nx.IsCashItem(v.ItemID) {
			p.WriteBytes(addItem(v, false))
		}
	}

	p.WriteByte(0)

	// Equips
	for _, v := range char.Equip {
		if v.SlotID < 0 && v.InvID == 1 && nx.IsCashItem(v.ItemID) {
			p.WriteBytes(addItem(v, false))
		}
	}

	p.WriteByte(0)

	// Inventory windows starts
	for _, v := range char.Equip {
		if v.SlotID > -1 && v.InvID == 1 {
			p.WriteBytes(addItem(v, false))
		}
	}

	p.WriteByte(0)

	for _, v := range char.Use {
		if v.InvID == 2 { // Use
			p.WriteBytes(addItem(v, false))
		}
	}

	p.WriteByte(0)

	for _, v := range char.SetUp {
		if v.InvID == 3 { // Set-up
			p.WriteBytes(addItem(v, false))
		}
	}

	p.WriteByte(0)

	for _, v := range char.Etc {
		if v.InvID == 4 { // Etc
			p.WriteBytes(addItem(v, false))
		}
	}

	p.WriteByte(0)

	for _, v := range char.Cash {
		if v.InvID == 5 { // Cash  - not working propery :(
			p.WriteBytes(addItem(v, false))
		}
	}

	p.WriteByte(0)

	// Skills
	p.WriteInt16(int16(len(char.Skills))) // number of skills

	for id, level := range char.Skills {
		p.WriteInt32(id)
		p.WriteInt32(level)
	}

	// Quests
	p.WriteInt16(0) // # of quests?

	// What are these for? Minigame record and some other things?
	p.WriteInt16(0)
	p.WriteInt32(0)
	p.WriteInt32(0)
	p.WriteInt32(0)
	p.WriteInt32(0)
	p.WriteInt32(0)

	p.WriteUint64(0)
	p.WriteUint64(0)
	p.WriteUint64(0)
	p.WriteUint64(0)
	p.WriteUint64(0)
	p.WriteInt64(0)

	return p
}

func addItem(item types.Item, shortSlot bool) maplepacket.Packet {
	p := maplepacket.NewPacket()

	if !shortSlot {
		if nx.IsCashItem(item.ItemID) && item.SlotID < 0 {
			p.WriteByte(byte(math.Abs(float64(item.SlotID + 100))))
		} else {
			p.WriteByte(byte(math.Abs(float64(item.SlotID))))
		}
	} else {
		p.WriteInt16(item.SlotID)
	}

	switch item.InvID {
	case 1:
		p.WriteByte(0x01)
	default:
		p.WriteByte(0x02)
	}

	p.WriteInt32(item.ItemID)

	if nx.IsCashItem(item.ItemID) {
		p.WriteByte(1)
		p.WriteUint64(uint64(item.ItemID))
	} else {
		p.WriteByte(0)
	}

	p.WriteUint64(item.ExpireTime)

	switch item.InvID {
	case 1:
		p.WriteByte(item.UpgradeSlots)
		p.WriteByte(item.ScrollLevel)
		p.WriteInt16(item.Str)
		p.WriteInt16(item.Dex)
		p.WriteInt16(item.Int)
		p.WriteInt16(item.Luk)
		p.WriteInt16(item.HP)
		p.WriteInt16(item.MP)
		p.WriteInt16(item.Watk)
		p.WriteInt16(item.Matk)
		p.WriteInt16(item.Wdef)
		p.WriteInt16(item.Mdef)
		p.WriteInt16(item.Accuracy)
		p.WriteInt16(item.Avoid)
		p.WriteInt16(item.Hands)
		p.WriteInt16(item.Speed)
		p.WriteInt16(item.Jump)
		p.WriteString(item.CreatorName)
		p.WriteInt16(item.Flag) // lock, show, spikes, cape, cold protection etc ?
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		p.WriteInt16(item.Amount)
		p.WriteString(item.CreatorName)
		p.WriteInt16(item.Flag) // lock, show, spikes, cape, cold protection etc ?
	default:
		fmt.Println("Unsuported item type", item.InvID)
	}

	return p
}

func writeDisplayCharacter(char types.Character) maplepacket.Packet {
	p := maplepacket.NewPacket()
	p.WriteByte(char.Gender) // gender
	p.WriteByte(char.Skin)   // skin
	p.WriteInt32(char.Face)  // face
	p.WriteByte(0x00)        // ?
	p.WriteInt32(char.Hair)  // hair

	cashWeapon := int32(0)

	for _, b := range char.Equip {
		if b.SlotID < 0 && b.SlotID > -20 {
			p.WriteByte(byte(math.Abs(float64(b.SlotID))))
			p.WriteInt32(b.ItemID)
		}
	}

	for _, b := range char.Equip {
		if b.SlotID < -100 {
			if b.SlotID == -111 {
				cashWeapon = b.ItemID
			} else {
				p.WriteByte(byte(math.Abs(float64(b.SlotID + 100))))
				p.WriteInt32(b.ItemID)
			}
		}
	}

	p.WriteByte(0xFF)
	p.WriteByte(0xFF)
	p.WriteInt32(cashWeapon)

	return p
}
