using System.Runtime.Serialization;

namespace famousfront.datamodels
{
internal class FeedMediaTypes
{
  public const uint Feed_media_type_none = 0;
  public const uint Feed_media_type_unknown = 1;
  public const uint Feed_media_type_url = 1 << 1;
  public const uint Feed_media_type_video = 1 << 2;
  public const uint Feed_media_type_audio = 1 << 3;
  public const uint Feed_media_type_image = 1 << 4;
  public const uint Feed_media_type_media = Feed_media_type_audio | Feed_media_type_video;
}
internal class FeedFlags
{
  public const uint Feed_flag_none = 0;
  public const uint Feed_flag_readed = 1;
  public const uint Feed_flag_star = 1 << 1;
  public const uint Feed_flag_save = 1 << 2;
}
	
internal class FeedStatuses
{
  public const ulong Feed_status_content_ready = 1ul << 0;                         //Feed_status_content_ready uint64 = 1 << iota
  public const ulong Feed_status_content_empty = 1ul << 1;                         //Feed_status_content_empty
  public const ulong Feed_status_content_inline = 1ul << 2;                        //Feed_status_content_inline
  public const ulong Feed_status_content_external_ready = 1ul << 3;                //Feed_status_content_external_ready
  public const ulong Feed_status_content_external_empty = 1ul << 4;                //Feed_status_content_external_empty
  public const ulong Feed_status_has_audio = 1ul << 5;                             //Feed_status_has_audio
  public const ulong Feed_status_has_video = 1ul << 6;                             //Feed_status_has_video
  public const ulong Feed_status_has_url = 1ul << 7;                               //Feed_status_has_url
  public const ulong Feed_status_has_image = 1ul << 8;                             //Feed_status_has_image
  public const ulong Feed_status_invisible = 1ul << 9;                             //Feed_status_invisible
  public const ulong Feed_status_text_empty = 1ul << 10;//deprecated               //Feed_status_text_empty //deprecated
  public const ulong Feed_status_text_little = 1ul << 11;                          //Feed_status_text_little
  public const ulong Feed_status_text_many = 1ul << 12;                            //Feed_status_text_many
  public const ulong Feed_status_image_empty = 1ul << 13;                          //Feed_status_image_empty
  public const ulong Feed_status_image_one = 1ul << 14;                            //Feed_status_image_one
  public const ulong Feed_status_image_many = 1ul << 15;                           //Feed_status_image_many
  public const ulong Feed_status_media_empty = 1ul << 16;// image, audio , video   Feed_status_media_empty // image, audio , video
  public const ulong Feed_status_media_one = 1ul << 17;                            //Feed_status_media_one
  public const ulong Feed_status_media_many = 1ul << 18;                           //Feed_status_media_many
  public const ulong Feed_status_media_inline = 1ul << 19;//deprecated             //Feed_status_media_inline //deprecated
  public const ulong Feed_status_linkdensity_low = 1ul << 20;                      //Feed_status_linkdensity_low
  public const ulong Feed_status_linkdensity_high = 1ul << 21;                     //Feed_status_linkdensity_high
  public const ulong Feed_status_format_flowdocument = 1ul << 22;                  //Feed_status_format_flowdocument
  public const ulong Feed_status_format_text = 1ul << 23;                          //Feed_status_format_text
  public const ulong Feed_status_mp4 = 1ul << 24;                                  //Feed_status_mp4
  public const ulong Feed_status_flv = 1ul << 25;                                  //Feed_status_flv
  public const ulong Feed_status_content_unresolved = 1ul << 26;                   //Feed_status_content_unresolved
  public const ulong Feed_status_summary_ready = 1ul << 27;                        //Feed_status_summary_ready
  public const ulong Feed_status_summary_empty = 1ul << 28;                        //Feed_status_summary_empty
  public const ulong Feed_status_content_unavail = 1ul << 29;                      //Feed_status_content_unavail
  public const ulong Feed_status_content_duplicated = 1ul << 30;                   //Feed_status_content_duplicated
	public const ulong Feed_status_summary_duplicated = 1ul << 31;                   //Feed_status_summary_duplicated
  public const ulong Feed_status_content_mediainline = 1ul << 32;                  //Feed_status_content_mediainline
  public const ulong Feed_status_summary_mediainline = 1ul << 33;                  //Feed_status_summary_mediainline
}
internal class FeedTypes
{
  public const uint Feed_type_unknown = 0;
  public const uint Feed_type_rss = 1 << 0;
  public const uint Feed_type_atom = 1 << 1;
  public const uint Feed_type_sinaweibo = 1 << 2;
  public const uint Feed_type_qqweibo = 1 << 3;
  public const uint Feed_type_blog = 1 << 4;
  public const uint Feed_type_tweet = 1 << 5;
  public const uint Feed_type_feed = Feed_type_rss | Feed_type_atom;
}
[DataContract]
internal class FeedLink
{
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string alias
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string cleaned_local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int words
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int sentences
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int links
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int density
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public long length
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public bool readable
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] images
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] videos
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] audios
  {
    get;set;
  }
}
[DataContract]
internal class FeedMedia
{
  [DataMember(EmitDefaultValue = false)]
  public string title
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string description
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int width
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int height
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int duration
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string mime
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string thumbanil
  {
    get;set;
  }
}
[DataContract]
internal class FeedAuthor
{
  [DataMember(EmitDefaultValue = false)]
  public string name
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string email
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong id
  {
    get;set;
  }
}
[DataContract]
internal class FeedTitle
{
  [DataMember(EmitDefaultValue = false)]
  public string main
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] others
  {
    get;set;
  }
}
[DataContract]
internal class FeedContent
{
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint words
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint density
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint links
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong status
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] images
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] medias
  {
    get;set;
  }
}
[DataContract]
internal class FeedEntry
{
  [DataMember(Name = "_id")]
  public string id
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint flags
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string source
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint type
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedTitle title
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedAuthor author
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public long pubdate
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string summary
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string content
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[]tags
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] images
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[]videos
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[]audios
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedLink[] links
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint words
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint density
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong status
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] categories
  {
    get;set;
  }
}
[DataContract]
internal class FeedSource
{
  [DataMember(EmitDefaultValue = false)]
  public string name
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint period
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong deadline
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint type
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public bool disabled
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public bool enable_proxy
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong update
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string website
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] tags
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] categories
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int unreaded
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia media
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string description
  {
    get;set;
  }
}
[DataContract]
internal class FeedTag
{
  [DataMember(Name = "value")]
  public string Name
  {
    get;set;
  }
}
internal class FeedSourceTypes
{
  public const string Feed_type_unknown = "";
  public const string Feed_type_rss = "rss";
  public const string Feed_type_atom = "atom";
  public const string Feed_type_blog = "blog";
  public const string Feed_type_tweet = "tweet";
  public const string Feed_type_sinaweibo = "weibo";
  public const string Feed_type_qqweibo = "qqweibo";
}
[DataContract]
internal class BackendError
{
    [DataMember]
    internal int code { get; set; }
    [DataMember]
    internal string reason { get; set; }
}
[DataContract]
internal class FeedsBackendConfig
{
    [DataMember(EmitDefaultValue = false)]
    public string ip { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public uint port { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string db_address { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string db_name { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string[] categories { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string data_dir { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public ulong usage { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string image { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string thumbnail { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string document { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string feed_source { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string feed_entry { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string proxy { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public uint summary_threshold { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public uint thumbnail_width { get; set; }
}
[DataContract]
internal class FeedsStatus
{
    [DataMember(EmitDefaultValue = false)]
    public ulong runned { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string error { get; set; }
}

[DataContract]
internal class FeedImage  // => ImageCache
{
    [DataMember(EmitDefaultValue = false)]
    public string mime { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string thumbnail { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string origin { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public int code { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public int width { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public int height { get; set; }
}
  [DataContract]
internal class FeedSourceFindEntry
{
    [DataMember(EmitDefaultValue = false)]
    public string url { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string title { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string summary { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string website { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public bool subscribed { get; set; }
}
}